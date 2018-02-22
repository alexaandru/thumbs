#!/usr/bin/env ruby

big, small = Dir['./content/images/articole/**/*.jpg'].sort.partition {|f| f !~ /_small..(jpg|png)/}
big.inject({}) {|s,v|
  small_name_template = File.join(File.dirname(v), File.basename(v, '.jpg') + '_small%s.jpg')
  small_names = [0,1].map{|i| small_name_template % i}
  small.include?(small_names.first) ? s : s.merge(v => small_names)
}.each {|src, dst|
  print "#{src} thumbnails... "
  [257, 566].each_with_index {|width, index|
    system(%|convert "#{src}" -resize #{width}x1000\\> -strip -quality 80 "#{dst[index]}"|)
  }
  puts "OK"
}
