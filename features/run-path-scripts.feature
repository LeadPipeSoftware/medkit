Feature: Create symlinks for files matching a convention in a dotfile repo
  As a user
  In order to manage my dotfiles from a central repo
  I want to create a symlink from a file in my dotfile repo

  Scenario: Symlink does not exist
    Given no symlink exists
     When I ask to create my symlinks
     Then symlinks are created

