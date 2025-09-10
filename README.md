# moldable 🔧

`moldable` builds precise interfaces from any package so you can plug in `mockery`, `gomock`, `moq`, or any other mock tool you like.

> [!WARNING]
> This project is in active development and may contain bugs or breaking changes. We recommend testing thoroughly in your environment before using in production. Issues and contributions are welcome!

## Why?

Most mock generators need an interface to do anything. When a library only exposes concrete structs you must hand-write that interface and keep it in sync with every upstream change. `moldable` creates it in one command and gets out of your way. After that you can use whatever mocking framework you like.
