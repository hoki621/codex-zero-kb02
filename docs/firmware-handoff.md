# 本人が進めるFirmware作業

zero-kb02の入力・表示を公開ライブラリで組み立て、USB CDC major 2でPC側Hostへ接続します。
Firmware本体は自分で実装します。PC側のコード・確認CLI・通信仕様・依存検証は用意されています。
完成済みの旧Firmwareはmajor 1です。新Hostの試験に流用しないでください。

## 1. PCと復旧手段を準備する（H0 / #27）

1. 親repoをsubmodule込みで取得し、`mise install`、`cd host && npm ci`を実行します。
2. `npm test`、`npm run device:check -- input`を実行します。USBを挿さなくても成功し、`MOCK`とキー1〜12、回転の累積`total=1`が出ることを確認します。
3. 親へ戻り、下記の依存検証を実行します。これは一時フォルダでビルドするだけです。

   ```sh
   mise exec -- sh docs/library-probe/check.sh
   ```

4. スクリプトが示す一時フォルダと`deps.txt`を確認します。HID mouseは含まれ、`machine/usb/hid/keyboard`は含まれません。**probe.uf2は書き込まないでください。製品Firmwareではありません。**
5. [既存の復旧手順](../firmware/docs/hardware-diagnostics.md)を読み、手元の現用UF2・復旧操作・配線写真を保存します。自分で実機作業を始めるまではflash不要です。

合格: mock成功、依存probeがビルドできる、復旧UF2と対象ボードを区別できる。

## 2. 入力ライブラリを準備する（C0成果 / H1の前提）

| 部品 | 固定版 | 呼び出すAPI / 任せる仕事 |
| --- | --- | --- |
| tinygo-keyboard | `cf173e98f60329b7f7feba941461bb95c065c418` + [input-only patch](patches/tinygo-keyboard-input-only.patch) | `New().AddMatrixKeyboard(...).Get()` / matrix走査とdebounce |
| drivers | v0.34.0 | `encoders.NewQuadratureViaInterrupt` / 回転検出、`ssd1306.NewI2C` / OLED |
| tinyfont / tinydraw | v0.6.0 / v0.4.0 | `WriteLine`、`Rectangle` / 文字と選択枠 |
| pio | v0.2.0 | `PIO0.ClaimStateMachine`、`piolib.NewWS2812B`、`WriteRaw` / LED波形 |
| TinyGo / Go | 0.40.1 / 1.25.13 | `machine.ADC`、`machine.Serial`、`machine/usb/hid/mouse` |

元のtinygo-keyboardはimportだけでVial用USB初期化を行い、`keyboard.go`から標準HID keyboardをimportします。
`New()`を呼ばないだけでは無効化できません。patchは`kb02_inputonly`タグで出力処理とVialを除外し、入力に必要な型だけを残します。
**matrixの走査・debounceファイルは一切変更していません。** 標準CDCHID descriptorには未使用のkeyboard項目が残りますが、keyboard handler・送信処理・Vial vendor interfaceを登録しません。USB descriptorの独自実装は追加しません。

1. `firmware/AGENTS.md`を「major 2・標準HID mouse許可・input-only利用」に合わせてから、自分のIssue branchで作業します。
2. 固定SHAをFirmware内の例 `third_party/tinygo-keyboard/` に取得し、上記patchを適用します。`LICENSE.txt`は必ず残します。
3. go.modに`replace github.com/sago35/tinygo-keyboard => ./third_party/tinygo-keyboard`を設定します。他の版とgo.sumは[検証module](library-probe/go.mod)を参照します。
4. buildには常に`-tags kb02_inputonly`を付けます。CI/buildコマンドも同じにします。タグなしでは元のUSB初期化が有効になるため、H1を進めないでください。
5. `Device.Tick()`や`Loop()`は使いません。input-onlyにこれらの出力dispatcherはありません。`MatrixKeyboard.Get()`が返す状態の変化を自分のKEYLOG/CDCに対応付けます。

`check.sh`は同じ固定SHA・patch・版でのfresh clone/buildを再現します。ライセンス・出典は[UPSTREAMS](../UPSTREAMS.md)。upstreamへの投稿はしていません。

## 3. 12キーだけを動かす（H1 / #28）

1. `AddMatrixKeyboard`へ列GP5/6/7/8、行GP9/10/11、keymap `nil`を渡します。実配線のダイオード方向が逆なら公開option `InvertDiode(true)`で合わせます。
2. まず約1ms周期で`Get()`を呼び、index=`row*4+column`を物理K1〜K12に対応付けます。debounceはライブラリが担います。自分で二重にdebounceしません。
3. `NoneToPress`で`KEYLOG <番号> DOWN`、`PressToRelease`でUPを改行付き出力します。押しっぱなしの`Press`から繰り返し送らないでください。この段階ではmajor 2を実装する必要はありません。
4. 自分のFirmwareをビルド・書込みした後、bridgeやserial monitorを閉じ、対象の完全なport名を確認します。
5. Host側から下記を実行し、各キーを20回押します。portの文字列は自分の機器に置き換えます。

   ```sh
   cd host
   npm run device:check -- raw --device --port /dev/cu.usbmodemYOUR_DEVICE
   ```

6. 押下・release各1回、番号が物理配置と一致、長押しで増えないことを記録します。Ctrl-Cで終了します。

合格: 12キー×20回で重複・欠落・番号違いなし。USB列挙にVial vendor interfaceがなく、キーを押してもmacOSへ文字が直接入力されないことも実機で確認します。ビルド成功だけで合格にしません。

## 4. 表示を作る（H2 / #29）

1. I2C0をSDA GP12/SCL GP13で設定し、SSD1306 128×64、address 0x3cで初期化します。上下が逆ならdriverのrotationで調整します。
2. `tinyfont.WriteLine`でW/I/B/D/U/E、`tinydraw.Rectangle`で選択枠を描きます。2列×3段の順は0/1、2/3、4/5です。変更がある時だけ`Display()`します。
3. PIO state machineを確保し、`NewWS2812B(..., GP1)`でLEDを設定します。`WriteRaw`はGRB配列です。色成分の上限を16/255にし、物理LED順を1灯ずつ確認します。FIFOへ12灯分を無条件に`PutRGB`して欠落させないでください。
4. 色はW青、I水色、B橙、D緑、U紫、E消灯とします。offline時はonlineと区別できる表示・消灯にします。
5. 受信部分をmajor 2のHELLO/STATE/OFFLINE/PINGまで実装したら、以下を実行します。

   ```sh
   npm run device:check -- display --device --port /dev/cu.usbmodemYOUR_DEVICE
   ```

合格: 6状態、選択枠0〜5、offlineが順番に見える。slot 5の選択確認時だけEをUへ変える仕様です。表示処理中もキーscanを止めないことを次の統合で確認します。

## 5. EncoderとJoystickを作る（H3 / #30）

1. EncoderをGP3/4で`NewQuadratureViaInterrupt`し、`Configure(QuadratureConfig{Precision:4})`から始めます。
2. `Position()`の前回との差を読みます。1クリックが1になるよう実機でPrecisionを合わせ、正が時計回りになるよう配線または符号を調整します。独自遷移表は作りません。
3. 差をmajor 2の`ENC <generation> <delta>`にします。1行は絶対値32以下、未送信量も有限にし、切断時は破棄します。
4. `machine.InitADC()`後、X=GP29、Y=GP28をConfigure/Getします。中立値を計測して保存し、dead zone、反転、速度上限を調整します。
5. 標準`mouse.Port().Move(dx,dy)`へ小さな相対移動量を渡します。両pushは未割当です。HID keyboardやVialを追加しません。
6. `input` CLIで低速/通常速度の左右各100クリックを数えます。

   ```sh
   npm run device:check -- input --device --port /dev/cu.usbmodemYOUR_DEVICE
   ```

合格: 実クリック数と累積量が一致、往復で元のtotalへ戻る、中立30秒でポインターが流れない。OLED更新を併用して再確認します。

## 6. major 2を統合する（H4 / #31）

1. [PROTOCOL.md](../PROTOCOL.md)を読み、128 byte行上限・整数範囲・handshake・generation・heartbeatを実装します。Go/TinyGo標準の文字列/数値/時間APIを使います。
2. 標準CDCの受信処理、matrix/Encoderの入力、表示更新を小さく分けます。表示は最新snapshotだけを保ち、入力処理に長い描画/通信待ちを入れません。
3. 起動・世代変更・再接続では押しっぱなしのキーを一度離すまで通知せず、古い回転を捨てます。offlineでは操作入力を送信しません。
4. `input`で通常入力を確認し、`faults`で不正行とheartbeat停止を確認します。

   ```sh
   npm run device:check -- faults --device --port /dev/cu.usbmodemYOUR_DEVICE
   ```

5. 初期STATEが表示され、過長行で暴走せず、13秒の無通信中に12秒を目安としてofflineになることを確認します。`input`を再起動して復帰させます。
6. `input`実行中にUSBを抜き差しします。自分で操作する時だけ行い、復帰時に押しっぱなしだったキーが新しい操作にならないことを確認します。

合格: protocolの正常/異常例を満たす、queue/行bufferが有限、通信断でoffline、再接続で全snapshotが復帰、表示中のキー/回転に欠落なし。

## 7. Herdr・Codexと統合する（X1 / #32）

1. 全確認CLIを終了します。READMEの順でbrew版専用server、Herdrの`codex-micro`、通常TerminalのHost bridgeを起動します。
2. K2/K3/K5〜8の6枠選択、K1 Escape、K4 popup、K12新規会話、Encoderを確認します。
3. 対応版0.155.1のcommand承認でK9/K10がそれぞれ1回だけ動くことを確認します。質問待ち・未対応承認・対象変更では送信されないことも確認します。K9は実際にcommandを許可するので、自分で内容を確認できる無害な試験だけにします。
4. USB/Host/Herdr再起動を分けて試し、別paneへ入力が飛ばないことを確認します。
5. Host/Firmware SHA、brew Codex版、Herdr版、UF2 SHA256、成功/失敗/未実施を#32へ記録します。発表デモを3回通し、復旧UF2と予備録画を保存します。

ここまでの実機試験は未実施です。#32の受入が終わるまで、製品全体の完成・実機互換性を主張しません。
