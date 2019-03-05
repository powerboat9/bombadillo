Control keys/input:

q     quit
j     scrolldown
k     scrollup
f     toggle showing favorites as subwindow
r     refresh current page data (re-request)

:#    go to link num
:url  go to url

:w # name       write linknum to file 
:w url name     write url to file
:w name         write current to file

:q              quit

:f add #__ name    add link num as favorite
:f add url name  add link url as favorite
:f add name      add current page as favorite
:f del #         delete favorite with num
:f del url       delete favorite with url
:f del name      delete favorite with name
:f #             visit favorite with num

:s ...kywds      search assigned engine with keywords

:home #          set homepage to link num
:home url        set homepage to url
:home            visit home

- - - - - - - - - - - - - - - - - - 

Config format:

[favorites]
colorfield.space ++ gopher://colorfield.space:70/
My phlog ++ gopher://circumlunar.space/1/~sloum/

[options]
homepage ++ gopher://sdf.org
searchengine ++ gopher://floodgap.place/v2/veronicasomething
savelocation ++ ~/Downloads/
httpbrowser ++ lynx
openhttp ++ true
