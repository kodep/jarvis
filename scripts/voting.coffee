# Description:
#   Повзоляет создавать голосование
#
# Commands:
#   hubot создай голосование <вариант1>, <вариант2>, ... - Создаст голосование
#   hubot закончи|закрой голосование - Закроет голосование
#   hubot открой голосование - Откроет голосование обратно
#   hubot покажи результаты голосования - Покажет результаты голосования
#   hubot покажи варианты голосования|Какие варианты голосования? - Покажет варианты голосования
#   hubot голосую|голос за <id варианта>|<название варианта2>
#
# Author:
#   khmm12@gmail.com
#

_ = require('lodash')

module.exports = (robot) ->
  voting = new Voting(robot)

  robot.respond /Создай голосование (.+)$/i, (msg) ->
    voting.createVoting(msg.match[1], msg)
  robot.respond /Закончи|Закрой голосование/i, (msg) ->
    voting.finishVoting(msg)
  robot.respond /Открой голосование/i, (msg) ->
    voting.openVoting(msg)
  robot.respond /Покажи результаты голосования/i, (msg) ->
    voting.sendResults(msg)
  robot.respond /(Покажи варианты голосования)|(Какие варианты голосования(\?)?)/i, (msg) ->
    voting.sendChoises(msg)
  robot.respond /(голосую|голос) за (.+)$/i, (msg) ->
    choice = msg.match[2].trim()
    voting.vote(choice, msg)

class Voting
  constructor: (@robot) ->
    @active = false
  createVoting: (rawChoices, msg) ->
    @choices = _.object(rawChoices.split(/, /).map((value, index) -> [index + 1, value.trim()]))
    ids = _.keys(@choices)
    @votes = _.object(_.map(ids, (value) -> [value, []]))
    @active = true
    msg.reply('Принято, мой господин')
  openVoting: (msg) ->
    return msg.reply('Сейчас нет голосований') unless @isCreated()
    return msg.reply('Голосвание и так открыто') if @isActive()
    @active = true
    msg.reply('Открыл, мой господин')
  finishVoting: (msg) ->
    return msg.reply('Сейчас нет голосований') unless @isCreated()
    return msg.reply('Голосвание и так закончено') unless @isActive()
    @active = false
    msg.reply('Принято, мой господин')
    @sendResults(msg)
  vote: (choice, msg) ->
    return msg.reply('Сейчас нет голосований') unless @isCreated()
    return msg.reply('Голосование уже закончено') unless @isActive()
    if /\d{1,2}$/i.test(choice)
      choiceID = parseInt(choice, 10)
    else
      choiceID = _.findKey(@choices, (_choice) -> choice is _choice)
    return msg.reply('Господин, такого варианта нет') unless @choices[choiceID]
    votersForChoice = @votes[choiceID]
    sender = @robot.brain.usersForFuzzyName(msg.message.user['name'])[0].name
    return msg.reply('Господин, Вы уже голосовали') if votersForChoice.indexOf(sender) isnt -1
    votersForChoice.push(sender)
    msg.reply('Ваш голос принят, Господин')

  isActive: -> @active
  isCreated: -> not _.isEmpty(@choices)
  doesAnybodyVoted: -> not _.isEmpty(_.values(@votes))
  sendChoises: (msg) ->
    return msg.reply('Сейчас нет голосований') unless @isCreated()
    response = _.map(@choices, (choice, id) -> "#{id}: #{choice}").join('\n')
    msg.send response
  sendResults: (msg) ->
    return msg.reply('Сейчас нет голосований') unless @isCreated()
    votingResults = @_getVotingResults()
    response = _.map(votingResults, (vote) ->
      "#{vote.id}: #{vote.choice} - #{vote.count}"
    ).join('\n')
    msg.send response

  _getVotingResults: ->
    # Transform votes object to [[id, [names]]]
    # then to [{id: <choice id>, choice: <name>, count: <votes number>, votes: [<nicknames of voters>]}]
    # sorted by votes number
    _(@votes).pairs()
    .map((pair) => { id: pair[0], choice: @choices[pair[0]], count: pair[1].length, votes: pair[1] })
    .sortBy((vote) -> -vote.count)
    .value()
