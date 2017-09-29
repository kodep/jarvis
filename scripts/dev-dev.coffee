# Description:
#   Код-код!
#
# Commands:
#   hubot код-код
#
# Author:
#   khmm12@gmail.com
#

GIFS = [
  'https://user-images.githubusercontent.com/4437249/31013381-7e89dc64-a51e-11e7-959a-84c9376e3e28.gif',
  'https://media.giphy.com/media/JIX9t2j0ZTN9S/giphy.gif',
  'https://media3.giphy.com/media/vhsNmFjuN4WDS/giphy.gif',
  'https://media3.giphy.com/media/13HBDT4QSTpveU/giphy.gif',
  'https://media2.giphy.com/media/LHZyixOnHwDDy/giphy.gif',
  'https://media2.giphy.com/media/11JTxkrmq4bGE0/giphy.gif',
  'https://www.emergingtechnologyadvisors.com/images/posts/2017-02-27-hour-of-code-2016/cat-typing.gif',
  'https://cs5.pikabu.ru/post_img/2015/09/14/9/1442242489_1321040046.gif',
  'http://javasea.ru/uploads/posts/2017-03/1490894193_kot-programmist.gif',
  'https://68.media.tumblr.com/d180debc05eb5c283927a04fc797db00/tumblr_op3374yFdY1r3lrt6o1_500.gif',
  'http://data.whicdn.com/images/72853592/original.gif',
  'https://media.giphy.com/media/3oKIPnAiaMCws8nOsE/giphy.gif',
  'http://stream1.gifsoup.com/view5/20140719/5073433/computer-and-cat-o.gif'
]

module.exports = (robot) ->
  robot.respond /код-код!*/i, (res) ->
    gif = GIFS[Math.floor(Math.random() * GIFS.length)]
    res.send(gif)
