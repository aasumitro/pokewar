# Changelog

All notable changes to this module will be documented in this file.
The format is based on Keep a [Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

VERSION MODEL: {NUMBER}XX(1).{NUMBER}XX(2).{NUMBER}XX(3)
- XX(1) -> Major Version
- XX(2) -> Minor Version
- XX(3) -> Patch Version

UPDATE VERSION:
- XX(1) -> Updated when breaking/make incompatible API changes
- XX(2) -> Updated when add functionality in a backwards compatible manner
- XX(3) -> Updated when patch and improvement

[0.0.3-dev]
---

- Fixing UI: 
  - Battle - action button section, hide when game started
- Fixing WS Handler:
  - save func - unlocked mutex can not be accessed
  - save func - battle_logs when game saved
  - sendMessage func - message_type replace with lib const instead of inject

[0.0.2-dev]
---

- UI Improvement
- WS Handler Improvement
- Add battle log after player annulled
- Update Test Coverage
- Fixing Player Rank & Score
- Fixing Monster Sync
