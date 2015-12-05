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
Promise = require('bluebird')

module.exports = (robot) ->
  boobsService = new Boobs()
  buttsService = new Butts()

  robot.respond /(покажи сиськи)|(show boobs)/i, (res) ->
    boobsService.getRandom()
    .then((result) -> res.send(result))
    .catch((err) ->
      robot.logger.error err
      res.reply('Я занят, попросите позже. Никаких сисек!')
    )

  robot.respond /(покажи попки)|(show butts)/i, (res) ->
    buttsService.getRandom()
    .then((result) -> res.send(result))
    .catch((err) ->
      robot.logger.error err
      res.reply('Я занят, попросите позже. Никаких попок!')
    )


class Pretty
  TTL: 5

  getRandom: ->
    @getCount()
    .then((count) => request.get("#{@API_URL}/#{_randomInt(count)}/1/rank"))
    .then((response) => "#{@CDN_URL}/#{response.data[0].preview}")

  getCount: ->
    time = new Date().valueOf()
    new Promise((resolve, reject) =>
      return resolve(@_count) if @_lastFetched? && (@_lastFetched + @TTL * 1000) - time > 0
      request.get("#{@API_URL}/count")
      .then((response) =>
        count = response.data[0].count
        resolve(count)
        @_count = count
        @_lastFetched = time
      )
      .catch((err) =>
         return resolve(@_count) if @_count?
         reject err
      )
    )

  _randomInt = (high) -> Math.floor(Math.random() * high)

class Boobs extends Pretty
  API_URL: 'http://api.oboobs.ru/boobs'
  CDN_URL: 'http://media.oboobs.ru'

class Butts extends Pretty
  API_URL: 'http://api.obutts.ru/butts'
  CDN_URL: 'http://media.obutts.ru'
