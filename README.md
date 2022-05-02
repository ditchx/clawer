Original post from
https://forum.golangbridge.org/t/golang-web-crawler/21311


Your job is to design and create a CLI application that will parser www.alexa.com  data. The application signature should look like the following:

$ ./clawer <action> <arg1> [<arg2>...] fg


The application must be able to accept these actions as param and perform the corresponding tasks:

    top : show top sites URL on www.alexa.com 
    country : show top 20 sites URL on www.alexa.comby country
    e.g.

$ ./clawer top
$ ./clawer country

The application needs to have an extensible interface where adding a new action is just a matter of adding more files and should
not require any modifications to the existing code base.
ps. If anything is unclear, you may set a reasonable assumption and state it at the beginning of the situation.


Since alexa.com has been retired, I used semrush.com instead.
