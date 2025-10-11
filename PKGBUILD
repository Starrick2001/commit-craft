# Maintainer: Your Name <youremail@domain.com>
pkgname=commit-craft
pkgver=1.2.0
pkgrel=1
pkgdesc="A tool to craft commit messages."
arch=('x86_64')
url="https://github.com/Starrick2001/commit-craft"
license=('MIT')
depends=('git')
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::$url/archive/refs/tags/v$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
  cd "$srcdir/$pkgname-$pkgver"
  go build -o "$pkgname"
}

package() {
  cd "$srcdir/$pkgname-$pkgver"
  install -Dm755 "$pkgname" "$pkgdir/usr/bin/$pkgname"
}
