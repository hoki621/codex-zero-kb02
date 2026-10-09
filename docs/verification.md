# 検証記録 / Verification — 2026-10-09

ビルド・mockの成功と、実機で観測した結果を分けています。 / Build and mock results are separate from physical observations.

## 実機で確認した版 / Physically tested revision

- Firmware: `9485033aac6c1e83098c19238cd1d31b98b9090b`
- Host: `1a9a8625d5bc11e4ea1df234e936703546a3ccf4`
- TinyGo 0.40.1 / Go 1.25.13
- UF2 SHA-256: `ee3f75332566795c20f858d7b6ef3436c6232ac0715c92ca6cd6e372d5cd4d27`
- USB serial: `zero-kb02-v2`; macOS port: `/dev/cu.usbmodemzero_kb02_v21`
- Herdr 0.9.3 / brew Codex CLI 0.162.0（撮影前確認 / pre-demo check）

上記UF2は元checkoutとfresh cloneで同じhashを確認してから、本人の許可を得て1回書き込みました。 / This UF2 matched between the original checkout and a fresh clone, and was flashed once with explicit permission.

| 確認 / Check | 結果 / Result |
| --- | --- |
| OLED / LED | 6枠・選択移動・offline・色と消灯を本人が目視確認 / Six slots, selection, offline and LEDs observed as expected |
| K1–K12 | 各20回、DOWN/UPが各20件で余分な通知なし / 20 presses each, exactly 20 DOWN and UP events per key |
| Encoder | CW/CCW各1クリックで期待イベント / One click each direction produced the expected event |
| Joystick | 約30秒中立でdriftなし・四方向正常（目視） / No drift over ~30 seconds and correct directions, visually observed |
| 過長行・heartbeat / Overlong line and heartbeat | 約12秒でoffline・LED消灯 / Offline and LEDs off after ~12 seconds |
| Herdr K2 / K4 / K12 | 1枠でfocus・popup開閉・新規会話を本人が画面確認 / Focus, popup and new chat observed in one trial slot |
| Encoder + Codex | effortが一段変わって戻り、modelは同じ / Effort changed and returned, model unchanged |
| USB再接続 / Reconnect | 単体とHost接続中に1回ずつ、表示とK2が復帰 / Display and K2 recovered after one reconnect in each mode |
| K9 / K10 | 0.162.0では未対応のため無効、未操作 / Disabled on unsupported 0.162.0; not exercised |

連続Encoder試験は手で約100クリックを数え、CW 127・CCW 11・net 116を記録しました。操作回数が不確かなため欠落率の証拠には使いません。ADC生値・表示負荷下の定量評価・6枠同時・K1実画面操作・押下中再接続・Host/Herdr再起動・live承認は未確認です。 / The approximate 100-click trial recorded CW 127, CCW 11, net 116; uncertain manual counting prevents a loss-rate conclusion. Raw ADC, display-load timing, six simultaneous slots, live K1, held-key reconnect, service restarts and live approval remain unverified.

## 公開整理版 / Publication cleanup

Firmwareの変更はDTR遷移の重複除去、定数・Go形式・理由コメントの整理です。Hostは出典コメントと文書の整理です。この整理版はビルドで確認し、実機へ再書き込みしていません。 / Firmware changes simplify DTR transitions and clarify constants, formatting and rationale. Host changes clarify provenance and docs. The cleanup is build-tested, not reflashed.

公開整理版の固定commit / Pinned cleanup commits:

- Firmware: `126c8a122831ae2d79b2c6ee9f89a131a7c62a3d`
- Host: `b3fa5086ce4e826e2e5482c79af6b77d59a85759`
- TinyGo 0.40.1 / Go 1.25.13 build: flash 40,576 bytes / RAM 14,588 bytes
- UF2 SHA-256: `6269fbe60a612a1cede0e4ec71231b931cda741c7462e764f510897fc545582d`
- Host typecheck・64 tests・mock dry-run、Firmware test/vet、input-only dependency check: PASS

検証コマンド / Checks:

```sh
mise exec -- sh -c 'cd host && npm run typecheck && npm test && npm run dry-run -- WIBDUE'
mise exec -- sh -c 'cd firmware && go test -count=1 ./... && go vet ./...'
mise exec -- sh -c 'cd firmware && tinygo build -o /tmp/zero-kb02.uf2 --target waveshare-rp2040-zero -tags kb02_inputonly --stack-size 8kb --size short .'
```

結果と最終SHAは[公開整理PR #37](https://github.com/hoki621/codex-zero-kb02/pull/37)に記録します。残る任意の実機確認は[#32](https://github.com/hoki621/codex-zero-kb02/issues/32)へ集約しています。 / Final SHAs and results are recorded in PR #37; optional remaining physical checks are consolidated in #32.

以前の隔離Codex API試験は0.155.1・0.160.0で成功しました。今回の0.162.0隔離smokeは実行環境で起動できず、実Herdr試験とは分けています。標準CDCHID descriptorには未使用keyboard項目が残りますが、input-onlyビルドはkeyboard handler・Vialを初期化しません。CDC writeは送達を保証しません。 / Earlier isolated API checks passed on 0.155.1 and 0.160.0. The 0.162.0 isolated smoke could not start in the execution environment. The standard CDCHID descriptor retains unused keyboard items, while the input-only build initializes neither a keyboard handler nor Vial. CDC writes do not guarantee delivery.
