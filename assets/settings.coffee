getColorsFromToneWells = ->
  tones = {}
  for n in [0..4]
    tones['tone'+n] = $('input[name="vsc-tone-'+n+'"]').val()
  return tones

toneColorsSaveClicked = ->
  saveColorsToCookies()

saveColorsToCookies = ->
  # console.log 'save colors'
  colors = getColorsFromToneWells()

  today = new Date()
  expire = new Date("2029-01-01 12:00:00")

  # set cookies
  for n in [0..4]
    document.cookie = 'tone'+n+'='+escape(colors['tone'+n])+'; expires='+expire.toGMTString()+"; path=/ "


$(document).ready( ->
  $('#tone-colors-save').on('click', -> toneColorsSaveClicked())

  # get tone colors from cookies or set default values
  for n in [0..4]
    if getCookie("tone"+n) == "undefined" or getCookie("tone"+n) == null
      switch n
        when 0 then $('input[name="vsc-tone-'+n+'"]').val("#000000")
        when 1 then $('input[name="vsc-tone-'+n+'"]').val("#00ac00")
        when 2 then $('input[name="vsc-tone-'+n+'"]').val("#021bff")
        when 3 then $('input[name="vsc-tone-'+n+'"]').val("#996633")
        when 4 then $('input[name="vsc-tone-'+n+'"]').val("#ff0000")
    else
      $('input[name="vsc-tone-'+n+'"]').val(unescape(getCookie("tone"+n)))

  saveColorsToCookies()
)


