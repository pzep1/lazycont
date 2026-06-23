class Lazycont < Formula
  desc "Lazydocker-style terminal UI for Apple's container CLI"
  homepage "https://github.com/pzep1/lazycont"
  url "https://github.com/pzep1/lazycont/archive/refs/tags/v0.2.0.tar.gz"
  sha256 "59aa85e1ecca349f031bb73cb02e7619a4c0bcd6f0cc3724d05444d9be351ca6"
  license "GPL-3.0-or-later"
  head "https://github.com/pzep1/lazycont.git", branch: "main"

  depends_on "go" => :build
  depends_on "container"
  depends_on :macos

  def install
    build_version = version.to_s
    build_version = "HEAD" if build_version.empty?

    system "go", "build",
      *std_go_args(ldflags: "-s -w -X main.version=#{build_version}"),
      "./cmd/lazycont"
  end

  def caveats
    <<~EOS
      lazycont drives Apple's container CLI. Homebrew installs it as a
      dependency, but you still need to start its system service before
      launching the TUI:

        brew services start container

      Or start it manually for the current session:

        container system start
    EOS
  end

  test do
    assert_match "lazycont", shell_output("#{bin}/lazycont --version")
    assert_match "Usage:", shell_output("#{bin}/lazycont --help")
  end
end
