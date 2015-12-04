ERROR_REPLIES = [ 'О, простите, я так устал', 'Кажется, что мне нужен отпуск',
            'Кажется, у меня байты не туда вставляются', 'Молодооой человееек, у нас обед', 'Я занят!' ]

MISS_REPLIES = [ 'Чего Вы хотите, мой господин?', 'Нипонимашки', 'Нипонимашки :(',
                 'Господин, может вам справку?', 'Хотел бы я Вам помочь, да я не понимаю' ]

module.exports = (robot) ->
  # https://github.com/github/hubot/issues/683
  # robot.respond /.+/, (res) ->
  #   res.reply res.random MISS_REPLIES
  robot.error (err, res) ->
    robot.logger.error err
    res.reply(res.random REPLIES) if res?
