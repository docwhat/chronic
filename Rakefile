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
    ' github.com/alecthomas/gometalinter'\
    ' github.com/mitchellh/gox'
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
  sh 'go build -v ./...'
end

desc 'Cross builds the executables'
task :xbuild do
  sh 'gox' \
  ' -os="linux windows"'\
  ' -arch="386 amd64 ppc64le"'\
  ' -osarch=darwin/amd64'\
  ' ./...'
end

desc 'Cleans up executables'
task :clean do
  rm FileList['chronic_*']
end
