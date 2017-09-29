# Description:
#   Код-код!
#
# Commands:
#   hubot код-код
#
# Author:
#   khmm12@gmail.com
#

module.exports = (robot) ->
  robot.send /код-код!*/i, (res) ->
    res.reply('https://media.giphy.com/media/JIX9t2j0ZTN9S/giphy.gif')
