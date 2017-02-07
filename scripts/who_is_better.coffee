# Description:
#   Позволяет узнать лучших сотрудников (в плане отработанного времени) за этот или за прошлый месяц
#
# Commands:
#   hubot кто молодец в этом месяце - Покажет Вам список сотрудников с их отработанным временем, по текущему месяцу
#   hubot кто молодец в прошлом месяце - Покажет список сотрудников с их отработанным временем, за прошлый месяц
#
# Author:
#   artem.telnov@kodep.ru
#

request = require('axios')
Promise = require('bluebird')

module.exports = (robot) ->
  whoIsBetter = new WhoIsBetterAPI()

  # Duration This Month
  #
  robot.respond /(кто молодец в этом месяце)/i, (res) ->
    whoIsBetter.getCount('duration', 'current')
    .then((result) -> res.send(result))
    .catch((err) ->
      robot.logger.error err
      res.reply 'Я занят, попросите позже'
    )

  # Duration Previous Month
  #
  robot.respond /(кто молодец в прошлом месяце)/i, (res) ->
    whoIsBetter.getCount('duration', 'previous')
    .then((result) -> res.send(result))
    .catch((err) ->
      robot.logger.error err
      res.reply 'Я занят, попросите позже'
    )


class WhoIsBetterAPI
  TTL: 5
  API_URL: "#{process.env.KODER_TIMER_API}/export_to_jarvis"
  getCount: (field, month) ->
    new Promise((resolve, reject) =>
      request.get("#{@API_URL}?field=#{field}&month=#{month}")
      .then((response) =>
        ArrayTheBest = response.data
        # поля с именем приводим к общему размеру. тем самым выравниваем поля с временем.
        ArrayTheBest = ArrayTheBest.map (i) ->  [ i[0] = i[0] + ' '.repeat(25 - i[0].length), i[1]].join ' '
        ArrayTheBest = ArrayTheBest.join '\n'
        @_best_name = ArrayTheBest
        resolve ArrayTheBest
      )
      .then((err) =>
         reject err
      )
    )
