require 'cocoapods-core'

podfile = Pod::Podfile.from_file('/Users/godrei/Develop/iOS/cards-up-ios/Podfile')
puts podfile.to_hash
puts podfile.workspace_path

