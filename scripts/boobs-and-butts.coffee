# Author:
#   khmm12@gmail.com
#

request = require('axios')

module.exports = (robot) ->
  boobsService = new Boobs()
  buttsService = new Butts()

  # Boobs
  #
  robot.respond /(покажи сиськи)|(хочу сисек)|(сиськи в студию)|(show boobs)/i, (res) ->
    boobsService.getRandom()
    .then((result) -> res.send(result))
    .catch((err) ->
      robot.logger.error err
      res.reply 'Я занят, попросите позже. Никаких сисек!'
    )

  robot.router.post '/hubot/boobs', (req, res) ->
    boobsService.getRandom()
    .then((result) -> res.json({text: result}))
    .catch((err) ->
      robot.logger.error err
      res.status(500).json({text: 'Я занят, попросите позже. Никаких сисек!'})
    )

  # Butts
  #
  robot.respond /(покажи попки)|(хочу попок)|(попки в студию)|(show butts)/i, (res) ->
    buttsService.getRandom()
    .then((result) -> res.send(result))
    .catch((err) ->
      robot.logger.error err
      res.reply 'Я занят, попросите позже. Никаких попок!'
    )

  robot.router.post '/hubot/butts', (req, res) ->
    buttsService.getRandom()
    .then((result) -> res.json({text: result}))
    .catch((err) ->
      robot.logger.error err
      res.status(500).json({text: 'Я занят, попросите позже. Никаких попок!'})
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
        @_count = count
        @_lastFetched = time
        resolve count
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
