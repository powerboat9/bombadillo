TODO
- Add styles/color support
- Add code comments/documentation for all items
- Make sure html links using the URL convention work correctly



Control keys/input:

q     quit
j     scrolldown
k     scrollup
f     toggle showing favorites as subwindow
TODO - r     refresh current page data (re-request)

GO
:#    go to link num
:url  go to url

SIMPLE
:quit                         quit
:home                         visit home
:bookmarks                    toogle bookmarks window
:search
:help

DOLINK
:delete #                     delete bookmark with num
:bookmarks #                  visit bookmark with num     

DOLINKAS
:write # name                 write linknum to file 
:add # name                   add link num as favorite

DOAS
:write url name               write url to file
:add url name                 add link url as favorite
:set something something      set a system variable



value, action, word

- - - - - - - - - - - - - - - - - - 

Config format:

[favorites]
colorfield.space ++ gopher://colorfield.space:70/
My phlog ++ gopher://circumlunar.space/1/~sloum/

[options]
home ++ gopher://sdf.org
searchengine ++ gopher://floodgap.place/v2/veronicasomething
savelocation ++ ~/Downloads/
httpbrowser ++ lynx
openhttp ++ true



