# 検証記録 — 2026-09-21

## 2026-10-09: major 2の実機確認

対象はFirmware `9485033aac6c1e83098c19238cd1d31b98b9090b`、Host `1a9a8625d5bc11e4ea1df234e936703546a3ccf4`。TinyGo 0.40.1 / Go 1.25.13で作ったUF2のSHA-256は`ee3f75332566795c20f858d7b6ef3436c6232ac0715c92ca6cd6e372d5cd4d27`。復旧用UF2を確認し、本人の許可を得て、接続中の旧`zero-kb02-v1`へ1回書き込みました。USBには`zero-kb02-v2`、`/dev/cu.usbmodemzero_kb02_v21`として再列挙しました。

| 確認 | 観測と判定 |
| --- | --- |
| 表示 | Hostの`display --device`を2回実行。6枠、選択枠の移動、offline、LEDの色と消灯を本人が目視で「すべて期待どおり」と確認。PASS（目視） |
| 入力 | Hostの`input --device`でK1〜K12を各20回。受信ログは各キーのDOWN/UPがそれぞれ20件で、番号違い・余分な通知なし。PASS（通信ログ） |
| Encoder | 1クリックずつのCW/CCWは各1件を受信。連続回転では本人がおよそ100クリックを手で数え、ログはCW 127件、CCW 11件、差し引き116。操作回数に誤差があるため100クリックの定量合否は保留 |
| Joystick | 中立約30秒でポインターのドリフトなし、上下左右が期待方向と本人が目視確認。PASS（目視。ADC生値は未測定） |
| 通信異常 | Hostの`faults --device`で過長行と13秒の通信停止を実行。HELLO major 2を受信、本人が約12秒でoffline表示・LED消灯を確認。PASS（CLIと目視） |
| USB抜き差し | 確認CLI終了後に本人が1回実施。表示の復帰を目視し、Macで同じ`/dev/cu.usbmodemzero_kb02_v21`の再列挙を確認。PASS（起動復帰のみ） |

実Herdr/Codexへの操作、Host接続中の状態復元、押下中の再接続、全不正入力、デモ3回と録画は未実施です。連続回転は厳密なクリック数を保証できないため、欠落率の証拠に使いません。

## 2026-10-09: 撮影前のHerdr接続確認

上記と同じFirmware `9485033` / Host `1a9a862` / UF2で、実機Host bridgeとHerdr 0.9.3を接続しました。brew Codex CLIは確認時点で0.162.0です。Hostの64テストが成功し、実Herdrの試験用`codex-micro`ペイン1枠を認識しました。専用App Serverは起動できましたが、隔離環境のsmoke試験は起動できず、今回の実ペインの結果を下記に分けて記録します。

| 確認 | 結果 |
| --- | --- |
| K2 | 試験用Codexペインへ切替、OLEDの1枠目を選択。本人が目視でPASS |
| Encoder左右1クリック | 試験用会話の推論の強さが一段変わって戻り、モデルは同じ。本人が画面でPASS |
| K4 | Herdrの状態一覧を開き、再押下で閉じる。本人が画面でPASS |
| K12 | idleの試験用ペインで新規会話へ移る。本人が画面でPASS |
| Host接続中のUSB抜き差し | キーを離して1回抜き差しし、OLEDとK2のペイン切替が復帰。本人が目視でPASS。押下中再接続は未実施 |
| K9/K10 | Codex CLI 0.162.0は承認画面の検証対象外で、Hostの診断も無効と表示。実行せず、撮影デモから省く |

これは1枠の撮影前確認です。6枠同時、承認/拒否、Host/Herdr再起動、3回連続デモと録画の合格記録ではありません。

## 2026-10-09: 実機を使わない組合せ検証

一時ディレクトリへ親mainをクリーン取得し、[Firmware PR #1](https://github.com/hoki621/codex-zero-kb02-firmware/pull/1)の`9485033aac6c1e83098c19238cd1d31b98b9090b`と[Host PR #3](https://github.com/hoki621/codex-zero-kb02-host/pull/3)の`1a9a8625d5bc11e4ea1df234e936703546a3ccf4`を子repoにチェックアウトしました。親mainのgitlinkはまだこの組合せを固定していません。

| 確認 | 結果 |
| --- | --- |
| Host `npm ci --ignore-scripts`、64テスト、mock `dry-run`/`input`/`display`/`faults` | PASS。一時npm cacheのみ使用 |
| Firmware `go test -count=1 ./...`、`go vet ./...` | PASS。一時Go cacheのみ使用 |
| Firmware TinyGo 0.40.1 / Go 1.25.13製品UF2 build | PASS、flash 40,544 byte、RAM 14,588 byte。UF2 SHA-256は`ee3f75332566795c20f858d7b6ef3436c6232ac0715c92ca6cd6e372d5cd4d27` |
| TinyGo依存一覧 | 標準CDC/HID mouseを含み、HID keyboardを含まない |
| 固定ライブラリprobe `sh docs/library-probe/check.sh` | PASS。固定SHAへpatchを適用し、HID keyboardなしでcompile-only UF2を生成。SHA-256は`1df270db82307d19566d559906e45966f02d83790394dc902e0a4f5fb0edec51` |
| brew Codex CLI 0.160.0 | 一時CODEX_HOMEでApp Serverと推論設定の変更がPASS、model turn 0。承認キーの実画面操作は未実施 |

最初の一時cloneではPATHがHomebrewのTinyGo 0.42.0とGo 1.25.14を選んでいました。以前この出力を0.40.1の結果と誤記したため訂正しました。実機用にはインストール済みの固定版バイナリをパス指定して再ビルドし、元checkoutとfresh cloneで上記hashが一致しました。brew/miseの導入済みツールや設定は変更していません。以下は2026-09-21時点の旧Host・旧Firmwareの記録として保持します。

PC側実装とFirmware着手用成果物の検証です。製品全体の実機受入ではありません。
Hostは `aeb1be89540271730140874566b916dde3a71ec8`（[PR #2](https://github.com/hoki621/codex-zero-kb02-host/pull/2)）です。正確なcommitは親gitlinkで固定します。Firmware `4d8104c5b4f3925b37394b0ca2d14c39486fc9a1` は未変更のmajor 1で、新Host major 2とは非互換です。

## 実施した確認

| 層 | 結果 |
| --- | --- |
| Host typecheck・clean build・全テスト | PASS、64テスト。fresh remote clone + npm ciでも再現。build前にdistを削除 |
| fake Herdr / USB / App Server | PASS、実機・実socketへ送信せず実装を検証 |
| 開発CLI | mock input/display/raw/faults、明示port必須、不正port拒否をテスト |
| 診断・dry-run | major 2、brew実体・版・socket、STATE全snapshotを確認 |
| brew Codex隔離試験 | 0.155.1、ephemeral会話、turns=0、effort=null→medium、model=gpt-6-astra |
| patch適用 | 固定library SHAへgit apply --check成功、matrix本体はbyte単位で無変更 |
| ライブラリfresh clone/build | check.sh成功。TinyGo 0.40.1 / Go 1.25.13、flash 28,896 byte、RAM 13,248 byte |
| USB静的依存 | 標準CDC/HID mouseを含み、machine/usb/hid/keyboardを含まない |
| 独立レビュー | 別エージェント2名でレビュー→修正→再レビュー。未解消指摘なし |
| 実Herdr UI・実USB・flash・物理入力・表示 | **NOT RUN**。#28〜#32で本人が受入 |

Codex試験は既存設定・履歴を使わない一時CODEX_HOMEの専用プロセスで実施し、終了後に削除しました。
承認UIは公式0.155.1 snapshot/keymapとmockで検証しています。live y/n送信は行っていません。

依存probeのUF2 SHA256は `1df270db82307d19566d559906e45966f02d83790394dc902e0a4f5fb0edec51`。
これは検証出力の識別値で、製品用・復旧用UF2ではありません。ビルドパス等による差を製品不具合と扱わないでください。

## 独立レビューで直した不具合

| 指摘 | 修正 / 回帰確認 |
| --- | --- |
| serial open timeout後の遅いerrorでprocess終了 | listenerを維持、遅延errorテスト |
| 古いHerdr接続の失敗が新接続を切断 | 接続世代ごとに失敗を処理、stop/start競合テスト |
| 遅いstartup snapshotが新状態を上書き | snapshot sequence比較、イベント競合テスト |
| 別会話の完了で未回答承認を消す | 明示応答/解決/lifecycleだけで失効、複数要求テスト |
| 最終pane.read中に消えた登録を有効扱い | 取得後に両ファイル/PIDを再検査、削除raceテスト |
| dead-owner登録がresumeを妨げる | 排他claimで死んだownerだけ回収、live/同時claimテスト |
| focus変更で拒否した後も回転queueが残る | changeEffort=falseでもqueue失効、連続入力テスト |
| 承認保存失敗で古い証拠が残る | launcherとremote CLIを終了してPID証拠失効、実launcher/fake serverへの障害注入テスト |

独立再レビューは通信契約、patch適用、write timeout、送信queue上限も確認しました。
追加dependency、独自WebSocket、独自matrix/USB driverは導入していません。

## 制約と残る受入

- claim中の強制終了で`.json.lock`が残る場合は、owner processを確認して手動回復します。live ownerを自動削除しません。
- 最終確認とHerdr送信は別APIです。確認から送信までの変化を原子的に排除するAPIはありません。
- 標準CDCHID descriptor内の未使用keyboard項目は残ります。Vial/keyboard handlerがないことと、実USB列挙の確認を区別します。
- GPIO方向、Joystick方向、LED順とUSBの起動復帰は上記の実機観察で確認しました。Encoderの定量一致、ADC生値、表示負荷下の入力、Host接続中のUSB復旧は未確認です。
- 製品Firmwareのfresh buildは上記の固定版で成功しました。C4の実機受入済みgitlink・デモ3回/録画は#32後です。PC成果物の統合を実機受入済み統合とは呼びません。

## 再現コマンド

親repoで`mise.toml`指定版をPATHから利用できる状態で（導入済みなら再インストール不要）:

```sh
cd host
npm ci
npm run typecheck
npm test
npm run dry-run -- WIBDUE
npm run device:check -- doctor
npm run smoke:codex
cd ..
sh docs/library-probe/check.sh
```

実機の合格条件は[本人向け手順](firmware-handoff.md)を使用してください。
