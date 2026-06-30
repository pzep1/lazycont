class Lazycontainer < Formula
  desc "Lazydocker-style terminal UI for Apple's container CLI"
  homepage "https://github.com/pzep1/lazycontainer"
  url "https://github.com/pzep1/lazycontainer/archive/refs/tags/v0.5.1.tar.gz"
  sha256 "9f207ab5b0cffc0dedd146e44ef3f8546cdef3a19f6c35d40005e4367039b6c4"
  license "GPL-3.0-or-later"
  head "https://github.com/pzep1/lazycontainer.git", branch: "main"

  depends_on "go" => :build
  depends_on "container"
  depends_on :macos

  def install
    build_version = version.to_s
    build_version = "HEAD" if build_version.empty?

    system "go", "build",
      *std_go_args(ldflags: "-s -w -X main.version=#{build_version}"),
      "./cmd/lazycontainer"
  end

  def caveats
    <<~EOS
      lazycontainer drives Apple's container CLI. Homebrew installs it as a
      dependency, but you still need to start its system service before
      launching the TUI:

        brew services start container

      Or start it manually for the current session:

        container system start
    EOS
  end

  test do
    assert_match "lazycontainer", shell_output("#{bin}/lazycontainer --version")
    assert_match "Usage:", shell_output("#{bin}/lazycontainer --help")
  end
end
