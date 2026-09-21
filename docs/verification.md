# 検証記録 — 2026-09-21

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
- GPIO方向、Encoder分解能、Joystick校正、LED順、表示負荷、USB復旧は実機で確認します。
- C4の実機受入済みgitlink・製品Firmwareのfresh build・デモ3回/録画は#32後です。PC成果物の統合を実機受入済み統合とは呼びません。

## 再現コマンド

親repoでmise install後:

```sh
cd host
npm ci
npm run typecheck
npm test
npm run dry-run -- WIBDUE
npm run device:check -- doctor
npm run smoke:codex
cd ..
mise exec -- sh docs/library-probe/check.sh
```

実機の合格条件は[本人向け手順](firmware-handoff.md)を使用してください。
