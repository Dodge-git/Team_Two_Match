🟢 ДЕНЬ 1 — Дизайн и контракты
1️⃣ Спроектировать модели (без кода, на бумаге / в Notion)

Sport

Team

Player

Match

MatchStatus (enum)

Пропиши:

поля

индексы

foreign keys

ограничения

2️⃣ Определить HTTP API

Таблица эндпоинтов:

Sport
GET /sports

Team
POST   /teams
GET    /teams
GET    /teams/:id
PUT    /teams/:id
DELETE /teams/:id

Player
POST   /players
GET    /players?team_id=
PUT    /players/:id
DELETE /players/:id

Match
POST   /matches
GET    /matches
GET    /matches/:id
POST   /matches/:id/start
POST   /matches/:id/finish
POST   /matches/:id/cancel
GET    /matches/active

3️⃣ Согласовать Kafka контракты

Ты публикуешь:

match.started
{
  "match_id": 1,
  "home_team": "string",
  "away_team": "string",
  "sport": "string"
}

match.ended
{
  "match_id": 1,
  "home_team": "string",
  "away_team": "string",
  "final_score": "2-1"
}


Ты принимаешь:

match.goal
{
  "match_id": 1,
  "team_id": 5,
  "player_id": 12,
  "minute": 78
}