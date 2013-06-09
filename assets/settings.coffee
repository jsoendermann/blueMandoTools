getColors = ->
  tones = {}
  for n in [0..4]
    tones['tone'+n] = $('input[name="vsc-tone-'+n+'"]').val()
  return tones

toneColorsSaveClicked = ->
  colors = getColors()

  today = new Date()
  expire = new Date("2029-01-01 12:00:00")

  # set cookies
  for n in [0..4]
    document.cookie = 'tone'+n+'='+escape(colors['tone'+n])+';expires='+expire.toGMTString()

$(document).ready( ->
  $('#tone-colors-save').on('click', -> toneColorsSaveClicked())
)
