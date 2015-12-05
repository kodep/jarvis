# Description:
#   Повзоляет посмотреть сиськи и попки
#
# Commands:
#   hubot покажи сиськи (show boobs) - Покажет Вам традиционный символ плодородия
#   hubot покажи попки (show butts) - Покажет Вам женские ягодицы
#
# Author:
#   khmm12@gmail.com
#

request = require('axios')

randomInt = (high) -> Math.floor(Math.random() * high)

module.exports = (robot) ->
  robot.respond /(покажи сиськи)|(show boobs)/i, (res) ->
    boobsUrl = "http://api.oboobs.ru/boobs/#{randomInt(5000)}/1/rank"
    request.get(boobsUrl)
    .then((response) -> "http://media.oboobs.ru/#{response.data[0].preview}")
    .then((result) -> res.send(result))
    .catch((err) -> res.reply('Я занят, попросите позже. Никаких сисек!'))

  robot.respond /(покажи попки)|(show butts)/i, (res) ->
    boobsUrl = "http://api.obutts.ru/butts/#{randomInt(5000)}/1/rank"
    request.get(boobsUrl)
    .then((response) -> "http://media.obutts.ru/#{response.data[0].preview}")
    .then((result) -> res.send(result))
    .catch((err) -> res.reply('Я занят, попросите позже. Никаких попок!'))
