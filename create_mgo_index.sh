#!/usr/bin/env mongo WordPress
db.moments.createIndex( { time:-1 } )
db.moments.createIndex( { to_top_time: -1,time:-1 } )
db.moments.createIndex( { to_top_time: -1, comment_num: -1, time:-1 } )
db.comments.createIndex( {time:-1 } )
db.users.createIndex( { time:-1 } )
