request = require('axios')
Promise = require('bluebird')

module.exports = (robot) ->
  durationService = new Duration()


  # Duration This Month
  #
  robot.respond /(кто молодец в этом месяце)/i, (res) ->
    durationService.getCount('duration', 'current')
    .then((result) -> res.send(result))
    .catch((err) ->
      robot.logger.error err
      res.reply 'Я занят, попросите позже'
    )

  robot.router.post '/hubot/duration_this_month', (req, res) ->
    durationService.getCount('duration', 'current')
    .then((result) -> res.json({text: result}))
    .catch((err) ->
      robot.logger.error err
      res.status(500).json({text: 'Я занят, попросите позже'})
    )
  # Approved This Month
  #
  # robot.respond /(кто самый оплачиваемый)/i, (res) ->
  #   durationService.getCount('approved_duration', 'current')
  #   .then((result) -> res.send(result))
  #   .catch((err) ->
  #     robot.logger.error err
  #     res.reply 'Я занят, попросите позже'
  #   )
  #
  # robot.router.post '/hubot/approved_duration_this_month', (req, res) ->
  #   durationService.getCount('approved_duration', 'current')
  #   .then((result) -> res.json({text: result}))
  #   .catch((err) ->
  #     robot.logger.error err
  #     res.status(500).json({text: 'Я занят, попросите позже'})
  #   )

  # Duration Previous Month
  #
  robot.respond /(кто молодец в прошлом месяце)/i, (res) ->
    durationService.getCount('duration', 'previos')
    .then((result) -> res.send(result))
    .catch((err) ->
      robot.logger.error err
      res.reply 'Я занят, попросите позже'
    )

  robot.router.post '/hubot/duration_previous_month', (req, res) ->
    durationService.getCount('duration', 'previos')
    .then((result) -> res.json({text: result}))
    .catch((err) ->
      robot.logger.error err
      res.status(500).json({text: 'Я занят, попросите позже'})
    )

  # Approved Previous Month
  #
  # robot.respond /(кто был самым оплачиваемый)|(кому мы скажем спасибо за прошлый месяц)|(кто в прошлом месяце принес денег)/i, (res) ->
  #   durationService.getCount('approved_duration', 'previos')
  #   .then((result) -> res.send(result))
  #   .catch((err) ->
  #     robot.logger.error err
  #     res.reply 'Я занят, попросите позже'
  #   )
  #
  # robot.router.post '/hubot/approved_duration_previous_month', (req, res) ->
  #   durationService.getCount('approved_duration', 'previos')
  #   .then((result) -> res.json({text: result}))
  #   .catch((err) ->
  #     robot.logger.error err
  #     res.status(500).json({text: 'Я занят, попросите позже'})
  #   )



class Pretty
  TTL: 5

  getCount: (field, month) ->
    time = new Date().valueOf()
    new Promise((resolve, reject) =>
      return resolve(@_best_name) if @_lastFetched? && (@_lastFetched + @TTL * 1000) - time > 0
      request.get("#{@API_URL}?field=#{field}&month=#{month}")
      .then((response) =>
        best_name = response.data
        best_name = best_name.map (i) ->  [ i[0] = i[0] + ' '.repeat(25 - i[0].length), i[1]].join ' '
        best_name = best_name.join '\n'
        @_best_name = best_name
        @_lastFetched = time
        resolve best_name
      )
      .catch((err) =>
         return resolve(@_best_name) if @_best_name?
         reject err
      )
    )


class Duration extends Pretty
  API_URL: 'https://timer.kodep.ru/export_to_jarvis'
