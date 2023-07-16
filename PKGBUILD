# Maintainer: Vasilj Milošević <eboye@sbb.rs>

pkgname=bas-celik
pkgver=1.0.0
pkgrel=1
pkgdesc='Serbian ID Card reader'
url='https://github.com/ubavic/bas-celik'
license=('MIT')
arch=('x86_64')
makedepends=('go' 'git')
_commit=8fe92be8efcd9a40f4d709a49cb4b39cbd758a1b
source=(
  "$pkgname::git+$url.git#commit=$_commit"
)
sha512sums=('SKIP')
options=('!lto')

build() {
  cd "$pkgname"
  go build -o bas-celik main.go
}

package() {
  cd "$pkgname"
  install -Dm 755 "$pkgname" -t "$pkgdir/usr/bin"
  install -Dm 644 README.md -t "$pkgdir/usr/share/doc/$pkgname"
  install -Dm 644 LICENSE -t "$pkgdir/usr/share/licenses/$pkgname"
}

