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
_ = require('lodash')

module.exports = (robot) ->
  whoIsBetter = new WhoIsBetterAPI()

  # Duration This Month
  #
  robot.respond /(кто молодец в этом месяце)/i, (res) ->
    whoIsBetter.query('duration', 'current')
    .then((result) -> res.send(result))
    .catch((err) ->
      robot.logger.error err
      res.reply 'Я занят, попросите позже'
    )

  # Duration Previous Month
  #
  robot.respond /(кто молодец в прошлом месяце)/i, (res) ->
    whoIsBetter.query('duration', 'previous')
    .then((result) -> res.send(result))
    .catch((err) ->
      robot.logger.error err
      res.reply 'Я занят, попросите позже'
    )


class WhoIsBetterAPI
  API_URL = "#{process.env.KODER_TIMER_API}/export_to_jarvis"
  EMPLOYEE_NAME_WIDTH = 25
  query: (field, month) ->
    request.get("#{API_URL}?field=#{field}&month=#{month}")
    .then((response) =>
      bestEmployees = response.data
      text = bestEmployees
      .map((entity) ->
        [employee, time] = entity
        entityText = [_.padEnd(employee, EMPLOYEE_NAME_WIDTH), time].join(' ')
        entityText
      )
      .join('\n')
      "```#{text}```"
    )
