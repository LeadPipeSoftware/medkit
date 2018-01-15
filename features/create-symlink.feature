Feature: Create a symlink to a file
  As a user
  In order to manage my dotfiles
  I want to create a symlink from a file in my SCM folder

  Scenario: Symlink does not exist
    Given no symlink exists
     When I ask to create my symlinks
     Then symlinks are created

