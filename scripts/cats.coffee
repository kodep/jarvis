# Description:
#   Котики
#
# Commands:
#   hubot (котика в студию)|(гифку с котом)|(покажи кота)
#
# Author:
#   khmm12@gmail.com
#

request = require('axios')
Promise = require('bluebird')

API_ENDPOINT = 'http://api.giphy.com/v1/gifs/random?api_key=dc6zaTOxFJmzC&tag=cats+animal'

module.exports = (robot) ->
  fetcher = new CatsFetcher
  robot.respond /(котика в студию)|(гифку с котом)|(покажи кота)/i, (res) ->
    fetcher.getRandom()
    .then((result) -> res.send(result))
    .catch((err) ->
      robot.logger.error err
      res.reply 'Коты устали('
    )

  robot.router.post '/hubot/cats', (req, res) ->
    fetcher.getRandom()
    .then((result) -> res.json({ text: result }))
    .catch((err) ->
      robot.logger.error err
      res.status(500).end()
    )

class CatsFetcher
  getRandom: ->
    request.get(API_ENDPOINT)
    .then((response) ->
      data = response.data.data
      data.image_url
    )
