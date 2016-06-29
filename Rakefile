# frozen_string_literal: true
BINARY = 'chronic'.freeze
GO_SOURCE = Rake::FileList['**/*.go'].tap do |src|
  src.exclude %r{^vendor/}
  src.exclude do |f|
    `git ls-files '#{f}'`.empty?
  end
end

task default: :build

desc 'Clean things up'
task :clean do
  sh 'git clean -xfd'
end

desc 'Fetch all dependencies'
task :setup do
  sh 'go version'
  sh 'go get -u -v '\
    ' github.com/alecthomas/gometalinter'
  sh 'gometalinter --install --update'
end

desc 'Check lint and style'
task :lint do
  sh 'gometalinter --deadline=1m --disable=gotype ./...'
end

desc 'Test the code'
task :test do
  sh 'go test -v'
end

desc 'Build for the native platform'
task :build do
  sh 'go install -v'
end

def find_binary(name, os, arch)
  gopath = ENV['GOPATH']
  [
    File.join(gopath, 'bin', "#{os}_#{arch}", name),
    File.join(gopath, 'bin', "#{os}_#{arch}", "#{name}.exe"),
    File.join(gopath, 'bin', name),
    File.join(gopath, 'bin', "#{name}.exe")
  ].select { |b| File.executable? b }.compact.first
end

def xbuild(os, arch)
  mkdir_p 'build'
  sh "env GOOS=#{os} GOARCH=#{arch} go install"
  binary = find_binary(BINARY, os, arch)
  mv binary, "build/#{BINARY}-#{os}-#{arch}"
end

desc 'Build for all supported platforms'
task xbuild: %w(
  build:mac64
  build:lin32
  build:lin64
  build:win32
  build:win64
  build:ppc64
  build:ppc64le
)

namespace :build do
  task(:mac64)   { xbuild 'darwin',  'amd64' }

  task(:lin32)   { xbuild 'linux',   '386' }
  task(:lin64)   { xbuild 'linux',   'amd64' }

  task(:win32)   { xbuild 'windows', '386' }
  task(:win64)   { xbuild 'windows', 'amd64' }

  task(:ppc64)   { xbuild 'linux',   'ppc64' }
  task(:ppc64le) { xbuild 'linux',   'ppc64le' }
end
