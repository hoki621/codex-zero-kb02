# USB CDC protocol — major 2

この仕様はHostと本人が実装するTinyGo Firmwareの共通契約です。
Firmwareは入力番号を送り、Hostが操作の意味を決めます。Joystickは標準HID mouseで直接相対移動を送り、CDCには送りません。

**major 1と非互換です。既存のFirmware（`4d8104c`）はmajor 1のため、このHostとは接続できません。**
旧形式を自動解釈しません。[旧仕様](docs/protocol-v1.md)は旧Firmwareを調べるための記録です。

## フレームと数値

- ASCII `0x20..0x7e`、単一空白区切り、大文字、LF終端。受信はCRLFも許可。
- **終端を含め128 byte以下**。LFなら本文127、CRLFなら本文126 byte以下。
- 空行、不正文字、余分な空白・引数、未知command、不正値は行全体を捨てる。状態も通信生存時刻も更新しない。
- 過長行は次のLFまで捨て、次の行から復帰する。受信バッファを無限に増やさない。
- unsigned整数は`0`または先頭ゼロなしの十進数。符号・小数・指数を許可しない。
- `generation`: 1..18446744073709551615（uint64）。大小比較せず一致だけを使う。
- `sequence`: 0..4294967295（uint32）、上限の次は0。
- `key`: 1..12。`selected`: 0..5、または未選択`-`。
- `delta`: -32..-1または1..32。`+1`、`-0`、`01`は不正。
- USB CDCは115200 baud指定。Hostはユーザー指定の完全なportだけを排他openし、自動探索しない。

## 起動と切断

1. Deviceは起動時offline。Hostはopen後`HELLO HOST 2`を送る。
2. Deviceは対応majorなら`HELLO ZERO-KB02 2`を返す。Hostの応答待ちは1秒。
3. 不一致major/不正応答/timeoutならHostはcloseする。KEY/ENCも操作も受け付けない。
4. 成功後、Hostは新しいgenerationの全snapshot（STATEまたはOFFLINE）とPINGを送る。
5. Deviceは有効なSTATEを受けてonlineになる。HELLOだけでは入力を送らない。
6. USB再接続、Host再起動、slotのterminal割当変更、online/offline変更でHostはgenerationを更新する。状態文字や選択枠だけの変化では更新しない。
7. Deviceはgeneration変更時に未送信入力・回転量を破棄する。既に押されているキーは一度離すまでDOWNを送らない。起動時も同じ。古い押下を新しい操作として再送しない。
8. Hostは切断時に入力contextを失効させ、同じ指定portへ1秒後から再接続する。復帰時には全snapshotを送る。

## Host → Device

| 行 | 意味 |
| --- | --- |
| `HELLO HOST 2` | 新しいhandshake。Deviceは未送信入力を破棄しofflineへ戻る |
| `STATE <generation> <selected> <states>` | 6枠全体のsnapshot。statesは`WIBDUE`の各文字を6個 |
| `OFFLINE <generation>` | Herdrの状態を取得できない。入力送信を止め、offline表示 |
| `PING <sequence>` | 同じsequenceをPONGで返す |

`W` working、`I` idle、`B` blocked（質問待ちも含む）、`D` done、`U` unknown、`E` empty。
selectedが指す枠はEであってはならない。DeviceはそのようなSTATEを行全体として拒否する。
STATEは変更時と5秒ごとに全体を送る。OFFLINEでもheartbeatは続ける。
有効なHELLO/STATE/OFFLINE/PINGを最後に受けてから12秒でDeviceはofflineに戻り、入力queueを破棄する。
HostはPINGを最大1個だけ未応答にし、そのsequenceに一致するPONGが12秒以内に来なければ切断する（検査間隔最大250ms）。

## Device → Host

| 行 | 意味 |
| --- | --- |
| `HELLO ZERO-KB02 2` | handshake応答。接続中に再受信した場合HostはDevice再起動として再接続 |
| `KEY <generation> <key> DOWN` | debounce済みの新しい押下 |
| `KEY <generation> <key> UP` | release。Hostの操作は行わない |
| `ENC <generation> <delta>` | 正=時計回り、負=反時計回り。1 detent=1段階 |
| `PONG <sequence>` | PINGへの応答 |

Hostはonlineかつ最新generationのKEY/ENCだけを扱う。古い世代は捨てて最新snapshotを再送する。
同じキーの連続DOWNは最初だけ受け付け、UPで解除する。長押しリピートはない。
Encoderは最大32段階を処理待ちにでき、超過は理由を記録して破棄する。
focus/thread/USB contextが変わった場合、待機回転を破棄する。動作の自動再試行はしない。
Firmwareは通信待ちでscanを止めず、最新表示への集約と有限の入力queueを使う。queue超過時は古い入力を再送せずofflineへ戻り、再handshakeを要求する。

| 物理キー | Host動作 |
| --- | --- |
| 1 | focus中の割当済みCodexへEscape |
| 2, 3, 5, 6, 7, 8 | slot 0, 1, 2, 3, 4, 5へfocus |
| 4 | Herdr session共通Status popup |
| 9, 10 | 対応する単独のcommand承認を確認できたときだけ固定y/nを1回 |
| 11 | 未割当 |
| 12 | focus中のidle/done Codexへ固定/new |

Encoder push・Joystick pushは送信しない。FirmwareはEscape/y/n等をHID keyboardとして送信しない。

## 相互実装のテスト例

以下の行末にはLFを付けます。正常例:

```text
HELLO HOST 2
HELLO ZERO-KB02 2
STATE 7 0 WIBDUE
PING 0
PONG 0
KEY 7 2 DOWN
KEY 7 2 UP
ENC 7 -3
STATE 7 - EEEEEE
OFFLINE 8
```

Device側の拒否例: `HELLO HOST 1`（非互換）、`STATE 0 0 WIBDUE`、`STATE 7 6 WIBDUE`、
`STATE 7 5 WIBDUE`（E選択）、`STATE 7 - WIBDU`（5枠）、`STATE 7 - WIBDUX`、`OFFLINE 07`、`PING -1`。
Host側の正常/異常入力は[実行されるベクトル一覧](host/test/protocol-vectors.ts)を正とします。
128 byte境界、分割/連結、CRLF、不正文字、再接続、heartbeatは[通信テスト](host/test/cdc-usb.test.ts)で検証します。

```sh
cd host
npm test
npm run device:check -- input
npm run device:check -- display
npm run device:check -- faults
```

すべてmockが既定です。実機の段階別試験は[本人向け手順](docs/firmware-handoff.md)へ進んでください。
