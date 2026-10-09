# Firmwareの段階別作業・実機受入

Firmware major 2の実装候補は[子repo PR #1](https://github.com/hoki621/codex-zero-kb02-firmware/pull/1)にあります。
以下のH0〜H4は部品ごとの確認・実機受入の順序として使います。コード担当はCodexに変更されました。
PRのビルド成功を実機合格とは扱いません。書き込みやport openは対象操作の明示依頼が必要です。

zero-kb02の入力・表示を公開ライブラリで組み立て、USB CDC major 2でPC側Hostへ接続します。
Firmware本体は子repo PR #1に実装済みです。以下の実装項目はコードとの照合用で、再実装は不要です。各段階の実機「合格」項目を順に記録します。
完成済みの旧Firmwareはmajor 1です。新Hostの試験に流用しないでください。

## 1. PCと復旧手段を準備する（H0 / #27）

1. 親repoをsubmodule込みで取得し、必要なツール版が既に揃っていれば再インストールせず、`cd host && npm ci`を実行します。不足版がある場合だけ親の`mise.toml`を確認して導入します。
2. `npm test`、`npm run device:check -- input`を実行します。USBを挿さなくても成功し、`MOCK`とキー1〜12、回転の累積`total=1`が出ることを確認します。
3. 親へ戻り、下記の依存検証を実行します。これは一時フォルダでビルドするだけです。

   ```sh
   sh docs/library-probe/check.sh
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

1. 子repoの`AGENTS.md`がmajor 2・標準HID mouse許可・input-only利用を定めていることを確認します。
2. `third_party/tinygo-keyboard/`に固定SHAとpatch、`LICENSE.txt`があり、go.modのreplaceがそのディレクトリを指すことを確認します。
3. 製品UF2をビルドするときは、子repo READMEのコマンドで必ず`-tags kb02_inputonly`を指定します。タグなしのビルドは失敗する設計です。
4. `MatrixKeyboard.Get()`が物理状態を取得し、Firmwareがその変化をCDCへ対応付ける境界をコードで確認します。

`check.sh`は同じ固定SHA・patch・版でのfresh clone/buildを再現します。ライセンス・出典は[UPSTREAMS](../UPSTREAMS.md)。upstreamへの投稿はしていません。

## 3. 12キーだけを動かす（H1 / #28）

1. コード上の列GP5/6/7/8・行GP9/10/11と物理K1〜K12の対応を確認します。ダイオード方向は実機で判定し、必要なら公開option `InvertDiode(true)`で修正します。
2. matrix走査とdebounceをライブラリへ任せ、押下・解放の変化だけを`KEY`として送ることを確認します。長押しで繰り返し送らないことは実機で検証します。
3. 製品Firmwareを書き込んだ後、bridgeやserial monitorを閉じ、対象の完全なport名を確認します。
4. Host側から下記を実行し、各キーを20回押します。portの文字列は自分の機器に置き換えます。

   ```sh
   cd host
   npm run device:check -- input --device --port /dev/cu.usbmodemYOUR_DEVICE
   ```

5. 押下・release各1回、番号が物理配置と一致、長押しで増えないことを記録します。Ctrl-Cで終了します。

合格: 12キー×20回で重複・欠落・番号違いなし。USB列挙にVial vendor interfaceがなく、キーを押してもmacOSへ文字が直接入力されないことも実機で確認します。ビルド成功だけで合格にしません。

## 4. 表示を作る（H2 / #29）

1. コード上のI2C0 GP12/13、SSD1306 128×64、address 0x3c、2列×3段の配置を確認します。表示の上下・文字の読みやすさは実機で判定します。
2. `tinyfont.WriteLine`で状態文字、`tinydraw.Rectangle`で選択枠、PIOの`WriteRaw`でGP1のLEDを駆動することを確認します。LED順と明るさは実機で判定します。
3. W青、I水色、B橙、D緑、U紫、E消灯とoffline表示を、以下で順に確認します。

   ```sh
   npm run device:check -- display --device --port /dev/cu.usbmodemYOUR_DEVICE
   ```

合格: 6状態、選択枠0〜5、offlineが順番に見える。slot 5の選択確認時だけEをUへ変える仕様です。表示処理中もキーscanを止めないことを次の統合で確認します。

## 5. EncoderとJoystickを作る（H3 / #30）

1. コード上のEncoder GP3/4、Precision 4、`ENC`の有限キューを確認します。1クリックの量と回転方向は実機で判定し、必要なら校正値を修正します。
2. JoystickのADC GP29/28、中立値、dead zone、標準HID mouseへの移動量を確認します。中立ドリフトと方向は実機で判定します。両pushは未割当です。
3. `input` CLIで低速/通常速度の左右各100クリックを数えます。

   ```sh
   npm run device:check -- input --device --port /dev/cu.usbmodemYOUR_DEVICE
   ```

合格: 実クリック数と累積量が一致、往復で元のtotalへ戻る、中立30秒でポインターが流れない。OLED更新を併用して再確認します。

## 6. major 2を統合する（H4 / #31）

1. [PROTOCOL.md](../PROTOCOL.md)とFirmwareの128 byte行上限・整数範囲・handshake・generation・heartbeatを照合します。
2. 押しっぱなしのキーや古い回転を起動・世代変更・再接続時に捨て、offlineでは操作を送らないことを、コードと実機の両方で確認します。
3. `input`で通常入力を確認し、`faults`で過長行とheartbeat停止を確認します。その他の異常系はGoテストと個別の実機確認に分けます。

   ```sh
   npm run device:check -- faults --device --port /dev/cu.usbmodemYOUR_DEVICE
   ```

4. 初期STATEが表示され、過長行で暴走せず、13秒の無通信中に12秒を目安としてofflineになることを確認します。`input`を再起動して復帰させます。
5. `input`実行中にUSBを抜き差しします。自分で操作する時だけ行い、復帰時に押しっぱなしだったキーが新しい操作にならないことを確認します。

合格: protocolの正常/異常例を満たす、queue/行bufferが有限、通信断でoffline、再接続で全snapshotが復帰、表示中のキー/回転に欠落なし。

## 7. Herdr・Codexと統合する（X1 / #32）

1. 全確認CLIを終了します。READMEの順でbrew版専用server、Herdrの`codex-micro`、通常TerminalのHost bridgeを起動します。
2. K2/K3/K5〜8の6枠選択、K1 Escape、K4 popup、K12新規会話、Encoderを確認します。
3. 対応版0.155.1または0.160.0のcommand承認でK9/K10がそれぞれ1回だけ動くことを確認します。質問待ち・未対応承認・対象変更では送信されないことも確認します。K9は実際にcommandを許可するので、自分で内容を確認できる無害な試験だけにします。
4. USB/Host/Herdr再起動を分けて試し、別paneへ入力が飛ばないことを確認します。
5. Host/Firmware SHA、brew Codex版、Herdr版、UF2 SHA256、成功/失敗/未実施を#32へ記録します。発表デモを3回通し、復旧UF2と予備録画を保存します。

ここまでの実機試験は未実施です。#32の受入が終わるまで、製品全体の完成・実機互換性を主張しません。
