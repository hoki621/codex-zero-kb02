# 発表デモの準備と記録

この説明は [Firmware PR #1](https://github.com/hoki621/codex-zero-kb02-firmware/pull/1) の `9485033` と [Host PR #3](https://github.com/hoki621/codex-zero-kb02-host/pull/3) の `1a9a862` を対象にします。親mainのgitlinkはまだこの組合せではありません。実機受入は [#32](https://github.com/hoki621/codex-zero-kb02/issues/32) に記録します。

2026-10-09の撮影前確認ではHerdr 0.9.3とbrew Codex CLI 0.162.0を使用し、試験用1枠でK2・K4・K12・Encoder・Host接続中のUSB復帰を確認しました。[結果](verification.md#2026-10-09-撮影前のherdr接続確認)を参照してください。0.162.0の承認画面はHostの検証対象外なので、今回の撮影でK9/K10を押さないでください。

## 説明するコードの流れ

| 見せる操作 | 公開ライブラリ・標準APIの役割 | この製品で書いた接続部分 |
| --- | --- | --- |
| キーを押してCodex paneを選ぶ | sago35/tinygo-keyboardの`MatrixKeyboard.Get()`が走査・debounce | Firmware `main.go`/`input.go`が物理番号を`KEY`にし、Host `usb.ts`/`bindings.ts`が現在のterminal identityを再確認してpaneを操作 |
| Encoderを回して推論の強さを変える | TinyGo driversの`NewQuadratureViaInterrupt`が回転量を取得 | Firmware `queue.go`が有限の`ENC`を送り、Host `reasoning.ts`が対象会話のeffortだけを変更 |
| Codexの状態をOLED/LEDに表示する | SSD1306、TinyFont、TinyDraw、PIOが描画とLED波形を担当 | Host `state.ts`/`usb.ts`が6枠の`STATE`を送り、Firmware `protocol.go`/`view.go`/`render.go`が選択枠・文字・色へ対応付け |
| Joystickでポインターを動かす | TinyGo標準HID mouseが相対移動を送る | Firmware `input.go`が中立値・dead zone・方向を調整。HostへのJoystick操作は送らない |

参照する実装は各PRのcommitで固定します。workshopはAPIの参考資料で、コードはコピーしていません。キーによる承認・拒否は、Hostが単独のcommand承認要求と対象・画面を確認できる場合だけです。

## 今回の撮影順（Codex CLI 0.162.0）

1. 現在のHost/Firmware SHA、brew Codex版、Herdr版、UF2 SHA-256を[#32](https://github.com/hoki621/codex-zero-kb02/issues/32)に記録する。復旧UF2と予備録画を手元に用意する。未確認の機器へこの手順だけを根拠に書き込まない。
2. 試験用Codexペイン1枠を用意し、K2でそのpaneとOLEDの選択枠へ移動する。K4で状態一覧を開閉し、Encoderを左右へ1クリックずつ回して対象会話のeffortだけが変わって戻ることを見せる。必要ならidleの試験用ペインでK12の新規会話も見せる。
3. Codex CLI 0.162.0ではK9/K10の承認・拒否を省き、未実施として記録する。対応版へ変更して再検証する場合だけ、内容を読める無害なcommand承認を別々に用意する。
4. Joystickでポインターを動かし、押しボタンには機能がないことを説明する。USBを抜き差しした後、表示と割当が復帰することを見せる。復帰しなければ復旧手順に移り、その回は失敗として記録する。

## リハーサル結果

成功だけを数え、連続3回の記録を残します。失敗時は原因と修正後の再実施を[#32](https://github.com/hoki621/codex-zero-kb02/issues/32)に記録します。

| 回 | Host/Firmware SHA・UF2 SHA-256 | 6枠・入力・表示 | 承認/拒否 | 再接続 | 結果・録画 |
| --- | --- | --- | --- | --- | --- |
| 1 | 未記録 | NOT RUN | NOT RUN | NOT RUN | NOT RUN |
| 2 | 未記録 | NOT RUN | NOT RUN | NOT RUN | NOT RUN |
| 3 | 未記録 | NOT RUN | NOT RUN | NOT RUN | NOT RUN |
