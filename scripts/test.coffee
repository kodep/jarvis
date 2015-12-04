module.exports = (robot) ->
  robot.respond /Привет/i, (res) ->
    res.reply('И тебе привет!')
