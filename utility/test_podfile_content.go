package utility

const testPodfileContent = `platform :ios, '9.0'
workspace 'Workspace'
use_frameworks!

abstract_target 'Applications' do
  project 'Apps'
  target 'WhatsTheScore'
  target 'Marcadores'

  target 'WorkspaceKit' do
    project 'WorkspaceKit'
    pod 'BNRDeferred', '~> 3.0.0-beta.3'
    target 'WorkspaceKitTests' do
      inherit! :search_paths
    end
  end

  target 'WorkspaceAPIKit' do
    project 'WorkspaceAPIKit'
    pod 'Alamofire'
    target 'WorkspaceAPIKitTests' do
      inherit! :search_paths
    end
  end

  target 'WorkspaceUIKit' do
    project 'WorkspaceUIKit'
    pod 'Cartography'

    target 'WorkspaceUIKitTests' do
      inherit! :search_paths
      pod 'FBSnapshotTestCase', ' ~> 2.1', '>= 2.1.4'
    end
  end
end
`
