!(function (e) {
    var n = {}
    function t(r) {
      if (n[r]) return n[r].exports
      var o = (n[r] = { i: r, l: !1, exports: {} })
      return e[r].call(o.exports, o, o.exports, t), (o.l = !0), o.exports
    }
    ;(t.m = e),
      (t.c = n),
      (t.d = function (e, n, r) {
        t.o(e, n) || Object.defineProperty(e, n, { enumerable: !0, get: r })
      }),
      (t.r = function (e) {
        "undefined" != typeof Symbol &&
          Symbol.toStringTag &&
          Object.defineProperty(e, Symbol.toStringTag, { value: "Module" }),
          Object.defineProperty(e, "__esModule", { value: !0 })
      }),
      (t.t = function (e, n) {
        if ((1 & n && (e = t(e)), 8 & n)) return e
        if (4 & n && "object" == typeof e && e && e.__esModule) return e
        var r = Object.create(null)
        if (
          (t.r(r),
          Object.defineProperty(r, "default", { enumerable: !0, value: e }),
          2 & n && "string" != typeof e)
        )
          for (var o in e)
            t.d(
              r,
              o,
              function (n) {
                return e[n]
              }.bind(null, o),
            )
        return r
      }),
      (t.n = function (e) {
        var n =
          e && e.__esModule
            ? function () {
                return e.default
              }
            : function () {
                return e
              }
        return t.d(n, "a", n), n
      }),
      (t.o = function (e, n) {
        return Object.prototype.hasOwnProperty.call(e, n)
      }),
      (t.p = "./webpack_dist/"),
      t((t.s = 7))
  })([
    function (e, n, t) {
      "use strict"
      n.a = function (e) {
        var n = this.constructor
        return this.then(
          function (t) {
            return n.resolve(e()).then(function () {
              return t
            })
          },
          function (t) {
            return n.resolve(e()).then(function () {
              return n.reject(t)
            })
          },
        )
      }
    },
    function (e, n) {
      var t
      t = (function () {
        return this
      })()
      try {
        t = t || Function("return this")() || (0, eval)("this")
      } catch (e) {
        "object" == typeof window && (t = window)
      }
      e.exports = t
    },
    function (e, n, t) {
      "use strict"
      ;(function (e) {
        var r = t(0),
          o = setTimeout
        function i() {}
        function a(e) {
          if (!(this instanceof a))
            throw new TypeError("Promises must be constructed via new")
          if ("function" != typeof e) throw new TypeError("not a function")
          ;(this._state = 0),
            (this._handled = !1),
            (this._value = void 0),
            (this._deferreds = []),
            d(e, this)
        }
        function c(e, n) {
          for (; 3 === e._state; ) e = e._value
          0 !== e._state
            ? ((e._handled = !0),
              a._immediateFn(function () {
                var t = 1 === e._state ? n.onFulfilled : n.onRejected
                if (null !== t) {
                  var r
                  try {
                    r = t(e._value)
                  } catch (e) {
                    return void u(n.promise, e)
                  }
                  s(n.promise, r)
                } else (1 === e._state ? s : u)(n.promise, e._value)
              }))
            : e._deferreds.push(n)
        }
        function s(e, n) {
          try {
            if (n === e)
              throw new TypeError("A promise cannot be resolved with itself.")
            if (n && ("object" == typeof n || "function" == typeof n)) {
              var t = n.then
              if (n instanceof a) return (e._state = 3), (e._value = n), void f(e)
              if ("function" == typeof t)
                return void d(
                  ((r = t),
                  (o = n),
                  function () {
                    r.apply(o, arguments)
                  }),
                  e,
                )
            }
            ;(e._state = 1), (e._value = n), f(e)
          } catch (n) {
            u(e, n)
          }
          var r, o
        }
        function u(e, n) {
          ;(e._state = 2), (e._value = n), f(e)
        }
        function f(e) {
          2 === e._state &&
            0 === e._deferreds.length &&
            a._immediateFn(function () {
              e._handled || a._unhandledRejectionFn(e._value)
            })
          for (var n = 0, t = e._deferreds.length; n < t; n++)
            c(e, e._deferreds[n])
          e._deferreds = null
        }
        function l(e, n, t) {
          ;(this.onFulfilled = "function" == typeof e ? e : null),
            (this.onRejected = "function" == typeof n ? n : null),
            (this.promise = t)
        }
        function d(e, n) {
          var t = !1
          try {
            e(
              function (e) {
                t || ((t = !0), s(n, e))
              },
              function (e) {
                t || ((t = !0), u(n, e))
              },
            )
          } catch (e) {
            if (t) return
            ;(t = !0), u(n, e)
          }
        }
        ;(a.prototype.catch = function (e) {
          return this.then(null, e)
        }),
          (a.prototype.then = function (e, n) {
            var t = new this.constructor(i)
            return c(this, new l(e, n, t)), t
          }),
          (a.prototype.finally = r.a),
          (a.all = function (e) {
            return new a(function (n, t) {
              if (!e || void 0 === e.length)
                throw new TypeError("Promise.all accepts an array")
              var r = Array.prototype.slice.call(e)
              if (0 === r.length) return n([])
              var o = r.length
              function i(e, a) {
                try {
                  if (a && ("object" == typeof a || "function" == typeof a)) {
                    var c = a.then
                    if ("function" == typeof c)
                      return void c.call(
                        a,
                        function (n) {
                          i(e, n)
                        },
                        t,
                      )
                  }
                  ;(r[e] = a), 0 == --o && n(r)
                } catch (e) {
                  t(e)
                }
              }
              for (var a = 0; a < r.length; a++) i(a, r[a])
            })
          }),
          (a.resolve = function (e) {
            return e && "object" == typeof e && e.constructor === a
              ? e
              : new a(function (n) {
                  n(e)
                })
          }),
          (a.reject = function (e) {
            return new a(function (n, t) {
              t(e)
            })
          }),
          (a.race = function (e) {
            return new a(function (n, t) {
              for (var r = 0, o = e.length; r < o; r++) e[r].then(n, t)
            })
          }),
          (a._immediateFn =
            ("function" == typeof e &&
              function (n) {
                e(n)
              }) ||
            function (e) {
              o(e, 0)
            }),
          (a._unhandledRejectionFn = function (e) {
            "undefined" != typeof console &&
              console &&
              console.warn("Possible Unhandled Promise Rejection:", e)
          }),
          (n.a = a)
      }).call(this, t(4).setImmediate)
    },
    function (e, n, t) {
      "use strict"
      ;(function (e) {
        var n = t(2),
          r = t(0),
          o = (function () {
            if ("undefined" != typeof self) return self
            if ("undefined" != typeof window) return window
            if (void 0 !== e) return e
            throw new Error("unable to locate global object")
          })()
        "Promise" in o
          ? o.Promise.prototype.finally || (o.Promise.prototype.finally = r.a)
          : (o.Promise = n.a)
      }).call(this, t(1))
    },
    function (e, n, t) {
      ;(function (e) {
        var r =
            (void 0 !== e && e) || ("undefined" != typeof self && self) || window,
          o = Function.prototype.apply
        function i(e, n) {
          ;(this._id = e), (this._clearFn = n)
        }
        ;(n.setTimeout = function () {
          return new i(o.call(setTimeout, r, arguments), clearTimeout)
        }),
          (n.setInterval = function () {
            return new i(o.call(setInterval, r, arguments), clearInterval)
          }),
          (n.clearTimeout = n.clearInterval =
            function (e) {
              e && e.close()
            }),
          (i.prototype.unref = i.prototype.ref = function () {}),
          (i.prototype.close = function () {
            this._clearFn.call(r, this._id)
          }),
          (n.enroll = function (e, n) {
            clearTimeout(e._idleTimeoutId), (e._idleTimeout = n)
          }),
          (n.unenroll = function (e) {
            clearTimeout(e._idleTimeoutId), (e._idleTimeout = -1)
          }),
          (n._unrefActive = n.active =
            function (e) {
              clearTimeout(e._idleTimeoutId)
              var n = e._idleTimeout
              n >= 0 &&
                (e._idleTimeoutId = setTimeout(function () {
                  e._onTimeout && e._onTimeout()
                }, n))
            }),
          t(5),
          (n.setImmediate =
            ("undefined" != typeof self && self.setImmediate) ||
            (void 0 !== e && e.setImmediate) ||
            (this && this.setImmediate)),
          (n.clearImmediate =
            ("undefined" != typeof self && self.clearImmediate) ||
            (void 0 !== e && e.clearImmediate) ||
            (this && this.clearImmediate))
      }).call(this, t(1))
    },
    function (e, n, t) {
      ;(function (e, n) {
        !(function (e, t) {
          "use strict"
          if (!e.setImmediate) {
            var r,
              o,
              i,
              a,
              c,
              s = 1,
              u = {},
              f = !1,
              l = e.document,
              d = Object.getPrototypeOf && Object.getPrototypeOf(e)
            ;(d = d && d.setTimeout ? d : e),
              "[object process]" === {}.toString.call(e.process)
                ? (r = function (e) {
                    n.nextTick(function () {
                      m(e)
                    })
                  })
                : !(function () {
                      if (e.postMessage && !e.importScripts) {
                        var n = !0,
                          t = e.onmessage
                        return (
                          (e.onmessage = function () {
                            n = !1
                          }),
                          e.postMessage("", "*"),
                          (e.onmessage = t),
                          n
                        )
                      }
                    })()
                  ? e.MessageChannel
                    ? (((i = new MessageChannel()).port1.onmessage = function (
                        e,
                      ) {
                        m(e.data)
                      }),
                      (r = function (e) {
                        i.port2.postMessage(e)
                      }))
                    : l && "onreadystatechange" in l.createElement("script")
                      ? ((o = l.documentElement),
                        (r = function (e) {
                          var n = l.createElement("script")
                          ;(n.onreadystatechange = function () {
                            m(e),
                              (n.onreadystatechange = null),
                              o.removeChild(n),
                              (n = null)
                          }),
                            o.appendChild(n)
                        }))
                      : (r = function (e) {
                          setTimeout(m, 0, e)
                        })
                  : ((a = "setImmediate$" + Math.random() + "$"),
                    (c = function (n) {
                      n.source === e &&
                        "string" == typeof n.data &&
                        0 === n.data.indexOf(a) &&
                        m(+n.data.slice(a.length))
                    }),
                    e.addEventListener
                      ? e.addEventListener("message", c, !1)
                      : e.attachEvent("onmessage", c),
                    (r = function (n) {
                      e.postMessage(a + n, "*")
                    })),
              (d.setImmediate = function (e) {
                "function" != typeof e && (e = new Function("" + e))
                for (
                  var n = new Array(arguments.length - 1), t = 0;
                  t < n.length;
                  t++
                )
                  n[t] = arguments[t + 1]
                var o = { callback: e, args: n }
                return (u[s] = o), r(s), s++
              }),
              (d.clearImmediate = p)
          }
          function p(e) {
            delete u[e]
          }
          function m(e) {
            if (f) setTimeout(m, 0, e)
            else {
              var n = u[e]
              if (n) {
                f = !0
                try {
                  !(function (e) {
                    var n = e.callback,
                      r = e.args
                    switch (r.length) {
                      case 0:
                        n()
                        break
                      case 1:
                        n(r[0])
                        break
                      case 2:
                        n(r[0], r[1])
                        break
                      case 3:
                        n(r[0], r[1], r[2])
                        break
                      default:
                        n.apply(t, r)
                    }
                  })(n)
                } finally {
                  p(e), (f = !1)
                }
              }
            }
          }
        })("undefined" == typeof self ? (void 0 === e ? this : e) : self)
      }).call(this, t(1), t(6))
    },
    function (e, n) {
      var t,
        r,
        o = (e.exports = {})
      function i() {
        throw new Error("setTimeout has not been defined")
      }
      function a() {
        throw new Error("clearTimeout has not been defined")
      }
      function c(e) {
        if (t === setTimeout) return setTimeout(e, 0)
        if ((t === i || !t) && setTimeout)
          return (t = setTimeout), setTimeout(e, 0)
        try {
          return t(e, 0)
        } catch (n) {
          try {
            return t.call(null, e, 0)
          } catch (n) {
            return t.call(this, e, 0)
          }
        }
      }
      !(function () {
        try {
          t = "function" == typeof setTimeout ? setTimeout : i
        } catch (e) {
          t = i
        }
        try {
          r = "function" == typeof clearTimeout ? clearTimeout : a
        } catch (e) {
          r = a
        }
      })()
      var s,
        u = [],
        f = !1,
        l = -1
      function d() {
        f &&
          s &&
          ((f = !1), s.length ? (u = s.concat(u)) : (l = -1), u.length && p())
      }
      function p() {
        if (!f) {
          var e = c(d)
          f = !0
          for (var n = u.length; n; ) {
            for (s = u, u = []; ++l < n; ) s && s[l].run()
            ;(l = -1), (n = u.length)
          }
          ;(s = null),
            (f = !1),
            (function (e) {
              if (r === clearTimeout) return clearTimeout(e)
              if ((r === a || !r) && clearTimeout)
                return (r = clearTimeout), clearTimeout(e)
              try {
                r(e)
              } catch (n) {
                try {
                  return r.call(null, e)
                } catch (n) {
                  return r.call(this, e)
                }
              }
            })(e)
        }
      }
      function m(e, n) {
        ;(this.fun = e), (this.array = n)
      }
      function h() {}
      ;(o.nextTick = function (e) {
        var n = new Array(arguments.length - 1)
        if (arguments.length > 1)
          for (var t = 1; t < arguments.length; t++) n[t - 1] = arguments[t]
        u.push(new m(e, n)), 1 !== u.length || f || c(p)
      }),
        (m.prototype.run = function () {
          this.fun.apply(null, this.array)
        }),
        (o.title = "browser"),
        (o.browser = !0),
        (o.env = {}),
        (o.argv = []),
        (o.version = ""),
        (o.versions = {}),
        (o.on = h),
        (o.addListener = h),
        (o.once = h),
        (o.off = h),
        (o.removeListener = h),
        (o.removeAllListeners = h),
        (o.emit = h),
        (o.prependListener = h),
        (o.prependOnceListener = h),
        (o.listeners = function (e) {
          return []
        }),
        (o.binding = function (e) {
          throw new Error("process.binding is not supported")
        }),
        (o.cwd = function () {
          return "/"
        }),
        (o.chdir = function (e) {
          throw new Error("process.chdir is not supported")
        }),
        (o.umask = function () {
          return 0
        })
    },
    function (e, n, t) {
      "use strict"
      t.r(n)
      t(3)
      function r(e) {
        var n = e.error,
          t = e.errormessage,
          r = e.scriptURI,
          o = e.initInfo,
          u = o.isIgnoreLimit,
          f = o.limitNum,
          l = o.onReport,
          d = o.ignoreError,
          p = o.versionName,
          m = o.userInfo,
          h = o.url,
          v = o.appName,
          y = o.packageName,
          g = o.svnVersion,
          w = o.storageUseName
        l && l(e)
        var T = window.localStorage.getItem(w)
            ? JSON.parse(window.localStorage.getItem(w))
            : { num: 0, timestamp: Date.now() },
          b = T.num
        if (
          !(
            !n ||
            (b > f && !u) ||
            !s(o) ||
            (function (e) {
              var n = e.error,
                t = void 0 === n ? {} : n,
                r = e.errormessage,
                o = void 0 === r ? "" : r,
                i = e.scriptURI,
                a = void 0 === i ? "" : i,
                c = e.ignoreError,
                s = !1
              Array.isArray(c) &&
                c.map(function (e) {
                  ;((t.stack && -1 !== t.stack.indexOf(e)) ||
                    (t.name && -1 !== t.name.indexOf(e)) ||
                    -1 !== o.indexOf(e) ||
                    -1 !== a.indexOf(e)) &&
                    (s = !0)
                })
              return s
            })({ error: n, errormessage: t, scriptURI: r, ignoreError: d })
          )
        ) {
          var S,
            O,
            I = {},
            E = i(),
            P = E.deviceType,
            x = E.OSVersion,
            _ = E.browserVersion,
            M = n.stack || n.message || ""
          ;(I.exceptionStack =
            "webViewVersion:" +
            _ +
            "\nurl:" +
            window.location.href +
            "\n" +
            M.slice(0, 3e4)),
            (I.exceptionType = n.name),
            (I.exceptionQualifier = c() + Math.floor(1e4 * Math.random())),
            (I.exceptionTime = c()),
            (I.networkType =
              navigator.connection && navigator.connection.type
                ? navigator.connection.type
                : ""),
            (I.cbasInfo = ""),
            (S =
              "userName=".concat(m.userName, "&") +
              "userID=".concat(
                (function (e) {
                  var n = ""
                  if (e) {
                    e = e.toString()
                    for (
                      var t = Math.floor(9 * Math.random()) + 1, r = 0;
                      r < e.length;
                      r++
                    ) {
                      var o = e.charAt(r).charCodeAt(),
                        i = String.fromCharCode(o + t)
                      n += i
                    }
                    n += t
                  }
                  return n
                })(m.userId),
                "&",
              ) +
              "systemType=".concat(a().toUpperCase(), "&") +
              "appName=".concat(v, "&") +
              "packageName=".concat(y, "&") +
              "versionName=".concat(p, "&") +
              "svnVersion=".concat(g, "&") +
              "phoneModel=".concat(P.trim(), "&") +
              "apiLevel=".concat(x.trim(), "&") +
              "monitorType=1&" +
              "monitorInfo=".concat(JSON.stringify(I))),
            (T.num = b + 1),
            window.localStorage.setItem(w, JSON.stringify(T)),
            ((O = { url: h, method: "POST", params: S }),
            new Promise(function (e, n) {
              ;(O.method = O.method ? O.method.toUpperCase() : "GET"),
                "GET" === O.method &&
                  (O.url +=
                    0 === location.search.length
                      ? "".concat("?", O.params)
                      : "".concat("&", O.params))
              var t = new XMLHttpRequest()
              ;(t.responseType = O.responseType || "json"),
                (t.onreadystatechange = function () {
                  if (4 === t.readyState)
                    return 200 === t.status ? e(t.response) : n(t.status)
                }),
                t.open(O.method, O.url, !0),
                "POST" === O.method &&
                  t.setRequestHeader(
                    "Content-Type",
                    "application/x-www-form-urlencoded",
                  ),
                t.send("POST" === O.method ? O.params : null)
            })).catch(function (e) {
              console.log("HXMonitor request error ", e)
            })
        }
      }
      function o() {
        var e,
          n = (function (e) {
            var n = e || window.navigator.userAgent,
              t = "",
              r = "",
              o = "",
              i = n.indexOf("Opera") > -1,
              a = n.indexOf("compatible") > -1 && n.indexOf("MSIE") > -1,
              c = n.indexOf("Trident") > -1 && n.indexOf("rv:11.0") > -1,
              s = n.indexOf("Edge") > -1,
              u = n.indexOf("Firefox") > -1,
              f = n.indexOf("360SE") > -1,
              l = !s && n.indexOf("Chrome") > -1 && n.indexOf("Safari") > -1,
              d = -1 === n.indexOf("Chrome") && n.indexOf("Safari") > -1
            if (
              (i
                ? ((o = "opera"), (t = "Opera"))
                : f
                  ? (t = "360SE")
                  : s
                    ? ((o = "edge"), (t = "Edge"))
                    : a
                      ? ((o = "MSIE"), (t = "IE"))
                      : c
                        ? ((o = "rv:"), (t = "IE11"))
                        : u
                          ? ((o = "Firefox"), (t = "Firefox"))
                          : d
                            ? ((o = "Safari"), (t = "Safari"))
                            : l
                              ? ((o = "Chrome"), (t = "Chrome"))
                              : (t = "Other"),
              o)
            ) {
              var p = new RegExp(o + "\\/?\\s*([\\d.]+)", "i"),
                m = n.match(p)
              null !== m && (r = m[1])
            }
            if (!r) {
              var h = n.match(/((rv\:)|(trident)|(msie)|(chrome))\/?\s*([\d.]+)/i)
              null !== h && (r = h[h.length - 1])
            }
            return { browserType: t, browserVersion: r }
          })()
        return {
          browserType: n.browserType,
          browserVersion: n.browserVersion,
          OSType:
            ((e = navigator.userAgent.toLowerCase()),
            "Win32" == navigator.platform || "Windows" == navigator.platform
              ? "Windows"
              : "Mac68K" == navigator.platform ||
                  "MacPPC" == navigator.platform ||
                  "Macintosh" == navigator.platform ||
                  "MacIntel" == navigator.platform
                ? "Mac"
                : e.indexOf("iPhone") > -1
                  ? "iPhone"
                  : e.indexOf("iPod") > -1
                    ? "iPod"
                    : e.indexOf("iPad") > -1
                      ? "iPad"
                      : e.indexOf("Android") > -1
                        ? "Android"
                        : "Other"),
        }
      }
      function i() {
        var e,
          n = (function () {
            var e = navigator.userAgent,
              n = a()
            if ("android" === n) {
              for (var t, r, o = e.split(";"), i = 0; i < o.length; i++)
                o[i].includes("Build/") && (t = i),
                  o[i].includes("Chrome/") && (r = i)
              var c = "",
                s = ""
              t &&
                ((c = (c = o[t]).slice(0, c.indexOf("Build/"))),
                (s = o[t - 1] ? o[t - 1] : ""))
              var u = ""
              return (
                r &&
                  (u = (u = o[r]).match(/Chrome\/(\S*)\s/)
                    ? u.match(/Chrome\/(\S*)\s/)[1]
                    : ""),
                { phoneModel: c, phoneOS: s, webViewVersion: u }
              )
            }
            return "iphone" === n
              ? {
                  phoneModel: "iPhone",
                  phoneOS: e.match(/iPhone\sOS\s(\S*)\s/)
                    ? e.match(/iPhone\sOS\s(\S*)\s/)[1]
                    : "",
                  webViewVersion: e.match(/AppleWebKit\/(\S*)\s/)
                    ? e.match(/AppleWebKit\/(\S*)\s/)[1]
                    : "",
                }
              : void 0
          })()
        if (n) {
          e = {
            deviceType: n.phoneModel,
            OSVersion: n.phoneOS,
            browserVersion: n.webViewVersion,
          }
        } else {
          var t = o(),
            r = t.browserType,
            i = t.browserVersion
          e = {
            deviceType: "PC",
            OSVersion: t.OSType,
            browserVersion: "".concat(r, "/").concat(i),
          }
        }
        return e
      }
      function a() {
        var e, n
        return (n = {
          iPhone:
            (e = navigator.userAgent).indexOf("iPhone") > -1 ||
            e.indexOf("Mac") > -1,
          iPad: e.indexOf("iPad") > -1,
          android: e.indexOf("Android") > -1,
        }).iPhone || n.iPad
          ? "iphone"
          : n.android
            ? "android"
            : "web"
      }
      function c() {
        function e(e) {
          return e < 10 ? "0" + e : e
        }
        var n = new Date()
        return (
          n.getFullYear().toString() +
          e(n.getMonth() + 1) +
          e(n.getDate()) +
          e(n.getHours()) +
          e(n.getMinutes()) +
          e(n.getSeconds())
        )
      }
      function s(e) {
        var n = e.versionName,
          t = e.userInfo,
          r = e.url,
          o = e.appName,
          i = e.packageName,
          a = e.svnVersion
        return n && t && t.userId && r && o && i && a
      }
      function u(e, n) {
        var t = null
        return function () {
          t
            ? Date.now() - t > n && ((t = Date.now()), e.apply(void 0, arguments))
            : ((t = Date.now()), e.apply(void 0, arguments))
        }
      }
      function f() {
        var e =
            navigator.connection && navigator.connection.type
              ? navigator.connection.type
              : "",
          n = i(),
          t = n.deviceType,
          r = n.OSVersion,
          o = n.browserVersion
        return [l(e), l(t), l(r), l(o)]
      }
      function l(e) {
        return e.replace(/(^\s+)|(\s+$)/g, "")
      }
      var d,
        p,
        m = function (e) {
          return {
            name: e,
            value:
              arguments.length > 1 && void 0 !== arguments[1] ? arguments[1] : -1,
            delta: 0,
            entries: [],
            id: ""
              .concat(Date.now(), "-")
              .concat(Math.floor(8999999999999 * Math.random()) + 1e12),
            isFinal: !1,
          }
        },
        h = function (e, n) {
          try {
            if (PerformanceObserver.supportedEntryTypes.includes(e)) {
              var t = new PerformanceObserver(function (e) {
                return e.getEntries().map(n)
              })
              return t.observe({ type: e, buffered: !0 }), t
            }
          } catch (e) {}
        },
        v = !1,
        y = !1,
        g = function (e) {
          v = !e.persisted
        },
        w = function (e) {
          var n = arguments.length > 1 && void 0 !== arguments[1] && arguments[1]
          y ||
            (addEventListener("pagehide", g),
            addEventListener("unload", function () {}),
            (y = !0)),
            addEventListener(
              "visibilitychange",
              function (n) {
                var t = n.timeStamp
                "hidden" === document.visibilityState &&
                  e({ timeStamp: t, isUnloading: v })
              },
              { capture: !0, once: n },
            )
        },
        T = function (e, n, t, r) {
          var o
          return function () {
            t && n.isFinal && t.disconnect(),
              n.value >= 0 &&
                (r || n.isFinal || "hidden" === document.visibilityState) &&
                ((n.delta = n.value - (o || 0)),
                (n.delta || n.isFinal || void 0 === o) && (e(n), (o = n.value)))
          }
        },
        b = function (e) {
          var n,
            t = arguments.length > 1 && void 0 !== arguments[1] && arguments[1],
            r = m("CLS", 0),
            o = function (e) {
              e.hadRecentInput || ((r.value += e.value), r.entries.push(e), n())
            },
            i = h("layout-shift", o)
          i &&
            ((n = T(e, r, i, t)),
            w(function (e) {
              var t = e.isUnloading
              i.takeRecords().map(o), t && (r.isFinal = !0), n()
            }))
        },
        S = function () {
          return (
            void 0 === d &&
              ((d = "hidden" === document.visibilityState ? 0 : 1 / 0),
              w(function (e) {
                var n = e.timeStamp
                return (d = n)
              }, !0)),
            {
              get timeStamp() {
                return d
              },
            }
          )
        },
        O = function (e) {
          var n,
            t = m("FCP"),
            r = S(),
            o = h("paint", function (e) {
              "first-contentful-paint" === e.name &&
                e.startTime < r.timeStamp &&
                ((t.value = e.startTime),
                (t.isFinal = !0),
                t.entries.push(e),
                n())
            })
          o && (n = T(e, t, o))
        },
        I = function (e) {
          var n = m("FID"),
            t = S(),
            r = function (e) {
              e.startTime < t.timeStamp &&
                ((n.value = e.processingStart - e.startTime),
                n.entries.push(e),
                (n.isFinal = !0),
                i())
            },
            o = h("first-input", r),
            i = T(e, n, o)
          o
            ? w(function () {
                o.takeRecords().map(r), o.disconnect()
              }, !0)
            : window.perfMetrics &&
              window.perfMetrics.onFirstInputDelay &&
              window.perfMetrics.onFirstInputDelay(function (e, r) {
                r.timeStamp < t.timeStamp &&
                  ((n.value = e),
                  (n.isFinal = !0),
                  (n.entries = [
                    {
                      entryType: "first-input",
                      name: r.type,
                      target: r.target,
                      cancelable: r.cancelable,
                      startTime: r.timeStamp,
                      processingStart: r.timeStamp + e,
                    },
                  ]),
                  i())
              })
        },
        E = function () {
          return (
            p ||
              (p = new Promise(function (e) {
                return ["scroll", "keydown", "pointerdown"].map(function (n) {
                  addEventListener(n, e, { once: !0, passive: !0, capture: !0 })
                })
              })),
            p
          )
        },
        P = function (e) {
          var n,
            t = arguments.length > 1 && void 0 !== arguments[1] && arguments[1],
            r = m("LCP"),
            o = S(),
            i = function (e) {
              var t = e.startTime
              t < o.timeStamp
                ? ((r.value = t), r.entries.push(e))
                : (r.isFinal = !0),
                n()
            },
            a = h("largest-contentful-paint", i)
          if (a) {
            n = T(e, r, a, t)
            var c = function () {
              r.isFinal || (a.takeRecords().map(i), (r.isFinal = !0), n())
            }
            E().then(c), w(c, !0)
          }
        },
        x = function (e) {
          var n,
            t = m("TTFB")
          ;(n = function () {
            try {
              var n =
                performance.getEntriesByType("navigation")[0] ||
                (function () {
                  var e = performance.timing,
                    n = { entryType: "navigation", startTime: 0 }
                  for (var t in e)
                    "navigationStart" !== t &&
                      "toJSON" !== t &&
                      (n[t] = Math.max(e[t] - e.navigationStart, 0))
                  return n
                })()
              ;(t.value = t.delta = n.responseStart),
                (t.entries = [n]),
                (t.isFinal = !0),
                e(t)
            } catch (e) {}
          }),
            "complete" === document.readyState
              ? setTimeout(n, 0)
              : addEventListener("pageshow", n)
        }
      function _(e) {
        for (var n = 1; n < arguments.length; n++) {
          var t = null != arguments[n] ? arguments[n] : {},
            r = Object.keys(t)
          "function" == typeof Object.getOwnPropertySymbols &&
            (r = r.concat(
              Object.getOwnPropertySymbols(t).filter(function (e) {
                return Object.getOwnPropertyDescriptor(t, e).enumerable
              }),
            )),
            r.forEach(function (n) {
              M(e, n, t[n])
            })
        }
        return e
      }
      function M(e, n, t) {
        return (
          n in e
            ? Object.defineProperty(e, n, {
                value: t,
                enumerable: !0,
                configurable: !0,
                writable: !0,
              })
            : (e[n] = t),
          e
        )
      }
      function F(e, n) {
        for (var t = 0; t < n.length; t++) {
          var r = n[t]
          ;(r.enumerable = r.enumerable || !1),
            (r.configurable = !0),
            "value" in r && (r.writable = !0),
            Object.defineProperty(e, r.key, r)
        }
      }
      window.Monitor = (function () {
        function e() {
          !(function (e, n) {
            if (!(e instanceof n))
              throw new TypeError("Cannot call a class as a function")
          })(this, e),
            (this.performanceException = [])
        }
        var n, t, o
        return (
          (n = e),
          (t = [
            {
              key: "install",
              value: function (e) {
                var n = this
                e &&
                  (e.config.errorHandler = u(function (e) {
                    r({ error: e, initInfo: n.initInfo })
                  }, 1e3))
              },
            },
            {
              key: "initMonitor",
              value: function (e) {
                var n = this
                ;(e.url =
                  e.url ||
                  "//mobileqa.hexin.cn/interface/MonitorInfoReceive.php"),
                  (this.initInfo = e)
                var t,
                  o = e.exceptionSwitch,
                  i = e.autoClearLocalFlag,
                  a = void 0 === i || i,
                  c = e.autoClearLocalDays,
                  s = void 0 === c ? 1 : c,
                  f = e.appName,
                  l = void 0 === f ? "" : f
                if (
                  ((this.initInfo.storageUseName = "".concat(
                    l,
                    ".collectExceptionNumData",
                  )),
                  o &&
                    ((window.onerror = u(function (e, t, o, i, a) {
                      r({
                        error: a,
                        errorMessage: e,
                        scriptURI: t,
                        initInfo: n.initInfo,
                      })
                    }, 1e3)),
                    window.addEventListener(
                      "error",
                      u(function (e) {
                        e.type &&
                          e.target.src &&
                          (n.performanceException.push({
                            type: e.type,
                            src: e.target.src,
                          }),
                          n.uploadStaticSourceError())
                      }, 1e3),
                      !0,
                    ),
                    window.addEventListener("unhandledrejection", function (e) {
                      try {
                        r({
                          error: {
                            name: e.reason.message,
                            stack: e.reason.stack,
                          },
                          initInfo: n.initInfo,
                        })
                      } catch (e) {}
                    })),
                  a && s)
                ) {
                  var d = window.localStorage.getItem(
                    this.initInfo.storageUseName,
                  )
                  d
                    ? ((d = JSON.parse(d)),
                      (t = d.timestamp),
                      Math.abs(t - Date.now()) / 864e5 > s &&
                        (d = { num: 0, timestamp: Date.now() }))
                    : (d = { num: 0, timestamp: Date.now() }),
                    window.localStorage.setItem(
                      this.initInfo.storageUseName,
                      JSON.stringify(d),
                    )
                }
              },
            },
            {
              key: "updateInitInfo",
              value: function (e) {
                e.url =
                  e.url || "//mobileqa.hexin.cn/interface/MonitorInfoReceive.php"
                var n = this.initInfo
                ;(this.initInfo = _({}, n, e)),
                  !s(n) && s(this.initInfo) && this.uploadStaticSourceError()
              },
            },
            {
              key: "collectExceptionData",
              value: function (e) {
                r(_({}, e, { initInfo: this.initInfo }))
              },
            },
            {
              key: "collectLoadDuration",
              value: function () {
                if (performance) {
                  var e = performance.timing,
                    n = e.loadEventEnd - e.navigationStart,
                    t = "normal",
                    r = [],
                    o = []
                  o.push(Math.round(n / 100)),
                    (o = o.concat(f())),
                    n > 2e3 && (t = "long")
                  var i = performance.getEntries()
                  if (i.length > 0) {
                    i.sort(function (e, n) {
                      return n.duration - e.duration
                    })
                    var a = window.location.href
                    i.forEach(function (e) {
                      if (
                        e.name !== a &&
                        /\.(js)|(css)$/.test(e.name) &&
                        e.duration > 100
                      ) {
                        var n = e.name.match(/\/([^\/]+)$/)
                        ;(n = n && n[1].length < 12 ? n[1] : e.name.slice(-12)),
                          r.push(
                            ""
                              .concat(n, ":")
                              .concat(Math.round(e.duration / 100)),
                          )
                      }
                    })
                  }
                  return { loadType: t, moreStr: r.join("|"), baseInfo: o }
                }
              },
            },
            {
              key: "collectWebVitals",
              value: function (e) {
                return new Promise(function (n, t) {
                  if (Array.isArray(e)) {
                    var r = [],
                      o = {
                        CLS: { poType: "layout-shift", func: b },
                        FCP: { poType: "paint", func: O },
                        FID: { poType: "first-input", func: I },
                        LCP: { poType: "largest-contentful-paint", func: P },
                        TTFB: { poType: "TTFB", func: x },
                      }
                    e.forEach(function (e) {
                      if (o[e]) {
                        var n = o[e],
                          t = n.poType,
                          i = n.func
                        !(function (e) {
                          if ("TTFB" === e) return !!window.performance
                          var n = !!navigator.userAgent.match(/Chrome|Chromium/),
                            t =
                              window.PerformanceObserver &&
                              Array.isArray(
                                PerformanceObserver.supportedEntryTypes,
                              ) &&
                              PerformanceObserver.supportedEntryTypes.includes(e)
                          return n && t
                        })(t)
                          ? r.push(Promise.resolve(""))
                          : r.push(
                              new Promise(function (e) {
                                i(function (n) {
                                  e(n)
                                })
                              }),
                            )
                      }
                    }),
                      Promise.all(r)
                        .then(function (e) {
                          n({ vitals: e, baseInfo: f() })
                        })
                        .catch(function (e) {
                          console.log("collectWebVitals catch e = ", e)
                        })
                  } else t("Parameter type error")
                })
              },
            },
            {
              key: "uploadStaticSourceError",
              value: function () {
                this.performanceException.length > 0 &&
                  s(this.initInfo) &&
                  (r({
                    initInfo: this.initInfo,
                    error: {
                      name: "sourceLoadError",
                      stack: JSON.stringify(this.performanceException),
                    },
                  }),
                  (this.performanceException = []))
              },
            },
          ]) && F(n.prototype, t),
          o && F(n, o),
          e
        )
      })()
    },
  ]) /*! 
  SkyWalking js-sdk
  @version 1.1.10
  @author wanghongwei,wangyifan
  @time 2023.03.22
   */
  !(function () {
    "use strict"
    var e,
      t,
      r = {
        "./src/index.ts": function () {
          function e() {
            return Math.random()
          }
          function t(e) {
            return Math.round(e)
          }
          function r(e) {
            return e >= 10 ? "".concat(e) : "0".concat(e)
          }
          function n() {
            var e,
              t = new Date()
            return {
              year: r(t.getFullYear()),
              month: r(t.getMonth() + 1),
              day: r(t.getDate()),
              hours: r(t.getHours()),
              minutes: r(t.getMinutes()),
              seconds: r(t.getSeconds()),
              miSeconds: ((e = t.getMilliseconds()), "0000".concat(e).slice(-3)),
            }
          }
          var o,
            i,
            a,
            c = {
              Chrome: /Chrome/,
              IE: /MSIE/,
              Firefox: /Firefox/,
              Opera: /Presto/,
              Safari: /Version\/([\d.]+).*Safari/,
              360: /360SE/,
              QQBrowser: /QQ/,
            },
            s = {
              iPhone: /iPhone/,
              iPad: /iPad/,
              Android: /Android/,
              Windows: /Windows/,
              Mac: /Macintosh/,
            }
          !(function (e) {
            ;(e.AJAX_ERROR = "ajax"),
              (e.RESOURCE_ERROR = "resource"),
              (e.VUE_ERROR = "vue"),
              (e.PROMISE_ERROR = "promise"),
              (e.JS_ERROR = "js"),
              (e.WHITE_SCREEN_ERROR = "whiteScreen"),
              (e.UNKNOWN_ERROR = "unknown")
          })(o || (o = {})),
            (function (e) {
              ;(e.INFO = "Info"), (e.WARNING = "Warning"), (e.ERROR = "Error")
            })(i || (i = {})),
            (function (e) {
              ;(e.ERROR = "/browser/errorELKLog"),
                (e.ERRORS = "/browser/errorELKLogs"),
                (e.PERF = "/browser/perfData"),
                (e.SEGMENT = "/v3/segment"),
                (e.SEGMENTS = "/v3/segments"),
                (e.PV = "/browser/sampleData")
            })(a || (a = {}))
          var u,
            l = "Http",
            f = "Exit"
          !(function (e) {
            ;(e[(e.OPENED = 1)] = "OPENED"), (e[(e.DONE = 4)] = "DONE")
          })(u || (u = {}))
          var p = 10001,
            d = "<browser>",
            h = navigator.userAgent,
            v = h.toLowerCase()
          function g() {
            var e,
              t = document.cookie
            if (0 !== t.length)
              for (var r = t.split(";"), n = 0; n < r.length; n++) {
                var o = r[n].split("="),
                  i = o[0],
                  a = o[1]
                "userid" === i.trim() && (e = a.trim())
              }
            return (
              e ||
              (function () {
                var e = v.indexOf("userid/")
                if (e > -1) return v.substr(e + 7).split(" ")[0]
              })()
            )
          }
          var m = function () {
              return (
                (m =
                  Object.assign ||
                  function (e) {
                    for (var t, r = 1, n = arguments.length; r < n; r++)
                      for (var o in (t = arguments[r]))
                        Object.prototype.hasOwnProperty.call(t, o) &&
                          (e[o] = t[o])
                    return e
                  }),
                m.apply(this, arguments)
              )
            },
            y = (function () {
              function e(e, t, r) {
                switch (
                  ((this.url = ""),
                  (this.type = ""),
                  (this.hook = void 0),
                  r && (this.hook = r),
                  (this.type = e),
                  e)
                ) {
                  case "ERROR":
                    this.url = t + a.ERROR
                    break
                  case "ERRORS":
                    this.url = t + a.ERRORS
                    break
                  case "SEGMENT":
                    this.url = t + a.SEGMENT
                    break
                  case "SEGMENTS":
                    this.url = t + a.SEGMENTS
                    break
                  case "PERF":
                    this.url = t + a.PERF
                    break
                  case "PV":
                    this.url = t + a.PV
                }
              }
              return (
                (e.prototype.sendByFetch = function (e) {
                  if ((delete e.collector, this.url)) {
                    var t = new Request(this.url, {
                      method: "POST",
                      body: JSON.stringify(e),
                    })
                    fetch(t)
                      .then(function (e) {
                        if (e.status >= 400 || 0 === e.status)
                          throw new Error("Something went wrong on api server!")
                      })
                      .catch(function (e) {
                        console.error(e)
                      })
                  }
                }),
                (e.prototype.sendByXhr = function (e) {
                  if (this.url) {
                    if (this.hook)
                      if (!this.hook({ type: this.type, url: this.url, data: e }))
                        return
                    var t = new XMLHttpRequest()
                    t.open("post", this.url, !0),
                      t.setRequestHeader("Content-Type", "application/json"),
                      (t.onreadystatechange = function () {
                        4 === t.readyState &&
                          t.status < 400 &&
                          console.log("Report successfully")
                      })
                    var r = ""
                    if (~["ERRORS", "ERROR"].indexOf(this.type)) {
                      var n = g(),
                        o = (function () {
                          var e = navigator.userAgent,
                            t = {
                              browser: "",
                              browserVersion: "",
                              os: "",
                              osVersion: "",
                              UA: e,
                            }
                          for (var r in c)
                            c[r].test(e) &&
                              ((t.browser = r),
                              "Chrome" === r
                                ? (t.browserVersion = e
                                    .split("Chrome/")[1]
                                    .split(" ")[0])
                                : "IE" === r
                                  ? (t.browserVersion = e
                                      .split("MSIE ")[1]
                                      .split(" ")[1])
                                  : "Firefox" === r
                                    ? (t.browserVersion = e.split("Firefox/")[1])
                                    : "Opera" === r
                                      ? (t.browserVersion =
                                          e.split("Version/")[1])
                                      : "Safari" === r
                                        ? (t.browserVersion = e
                                            .split("Version/")[1]
                                            .split(" ")[0])
                                        : "360" === r
                                          ? (t.browserVersion = "")
                                          : "QQBrowser" === r &&
                                            (t.browserVersion = e
                                              .split("Version/")[1]
                                              .split(" ")[0]))
                          for (var r in s)
                            s[r].test(e) &&
                              ((t.os = r),
                              "Windows" === r
                                ? (t.osVersion = e
                                    .split("Windows NT ")[1]
                                    .split(";")[0])
                                : "Mac" === r
                                  ? (t.osVersion = e
                                      .split("Mac OS X ")[1]
                                      .split(")")[0])
                                  : "iPhone" === r
                                    ? (t.osVersion = e
                                        .split("iPhone OS ")[1]
                                        .split(" ")[0])
                                    : "iPad" === r
                                      ? (t.osVersion = e
                                          .split("iPad; CPU OS ")[1]
                                          .split(" ")[0])
                                      : "Android" === r &&
                                        (t.osVersion = e
                                          .split("Android ")[1]
                                          .split(";")[0]))
                          return t
                        })(),
                        i = []
                      e.forEach(function (e) {
                        var t = e[0],
                          a = e[1]
                        if ("ajax" === t) i.push(a)
                        else if ("promise" === t && i.indexOf(a) > -1) return
                        if (
                          (n && (e.userid = n), (e = m(m({}, e), o)).pagePath)
                        ) {
                          var c = e.pagePath.replace(/(\d+)|./g, function (e, t) {
                              return t ? "*" : e
                            }),
                            s = e.pagePath.startsWith("/") ? "" : "/"
                          e.groupedPagePath = location.origin + s + c
                        }
                        r += "".concat(JSON.stringify(e), "\n")
                      })
                    } else r = JSON.stringify(e)
                    t.send(r)
                  }
                }),
                e
              )
            })(),
            w = new ((function () {
              function e() {
                ;(this.queues = []), (this.collector = ""), (this.hook = void 0)
              }
              return (
                (e.prototype.addTask = function (e, t, r) {
                  this.queues.push(e), (this.collector = t), (this.hook = r)
                }),
                (e.prototype.fireTasks = function () {
                  this.queues &&
                    this.queues.length &&
                    (new y("ERRORS", this.collector, this.hook).sendByXhr(
                      this.queues,
                    ),
                    (this.queues = []))
                }),
                (e.prototype.finallyFireTasks = function () {
                  var e = this
                  window.onbeforeunload = function () {
                    e.queues.length &&
                      new y("ERRORS", e.collector, e.hook).sendByXhr(e.queues)
                  }
                }),
                e
              )
            })())(),
            E = !1,
            O = (function () {
              function e() {
                this.logInfo = {
                  uniqueId: "",
                  category: o.UNKNOWN_ERROR,
                  service: "",
                  serviceVersion: "",
                  group: "",
                  subGroup: "",
                  pagePath: "",
                  grade: i.INFO,
                  errorUrl: "",
                  line: 0,
                  col: 0,
                  message: "",
                  firstReportedError: !1,
                  collector: "",
                  reportHook: void 0,
                }
              }
              return (
                (e.prototype.traceInfo = function (e) {
                  E ||
                    this.logInfo.category !== o.JS_ERROR ||
                    ((E = !0), (this.logInfo.firstReportedError = !0))
                  var t = this.logInfo,
                    r = t.collector,
                    n = t.reportHook
                  delete this.logInfo.collector,
                    w.addTask(this.logInfo, r, n),
                    w.finallyFireTasks(),
                    setTimeout(function () {
                      w.fireTasks()
                    }, e || 3e4)
                }),
                e
              )
            })()
          function b() {
            return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(
              /[xy]/g,
              function (t) {
                var r = (16 * e()) | 0
                return ("x" === t ? r : (3 & r) | 8).toString(16)
              },
            )
          }
          var R,
            S,
            _ = location.host,
            k = "v1.1.10",
            T = function (e, t) {
              void 0 === t && (t = [])
              var r,
                o = {
                  uniqueId: window.ELK_WHITESCREEN_TRACEID || b(),
                  collector: "",
                  service: "",
                  serviceVersion: "",
                  group: "",
                  subGroup: "",
                  pagePath: "",
                  domain: _,
                  sdkVersion: k,
                }
              ;((r = [
                "collector",
                "service",
                "serviceVersion",
                "group",
                "subGroup",
                "pagePath",
                "reportHook",
              ]),
              t.forEach(function (e) {
                ;-1 === r.indexOf(e) && r.push(e)
              }),
              r).forEach(function (t) {
                var r
                if ("pagePath" !== t || o[t])
                  o[t] =
                    "reportHook" === t
                      ? (null === (r = e.hooks) || void 0 === r
                          ? void 0
                          : r.reportHook) || e.reportHook
                      : e[t]
                else {
                  var n = location.pathname,
                    i = location.hash,
                    a = ""
                  n && "/" !== n && (a += n),
                    i && (a += n ? i : i.replace("#", "")),
                    (o[t] = a || "/")
                }
              })
              var i = n(),
                a = i.year,
                c = i.month,
                s = i.day,
                u = i.hours,
                l = i.minutes,
                f = i.seconds,
                p = i.miSeconds,
                d = ""
                  .concat(a, "-")
                  .concat(c, "-")
                  .concat(s, " ")
                  .concat(u, ":")
                  .concat(l, ":")
                  .concat(f, ".")
                  .concat(p)
              return (o.logtime = d), o
            },
            I = function (e, t) {
              var r = t,
                n = e.extraMsg
              return (
                n &&
                  ("function" == typeof n && (n = n()),
                  (r = "".concat(r, " -(").concat(n, ")"))),
                r
              )
            },
            P = function (e, t) {
              var r,
                n =
                  (null === (r = e.hooks) || void 0 === r
                    ? void 0
                    : r.filterHook) || e.filterHook
              if ("function" == typeof n)
                try {
                  return n(t)
                } catch (e) {
                  return !0
                }
              return !0
            },
            x =
              ((R = function (e, t) {
                return (
                  (R =
                    Object.setPrototypeOf ||
                    ({ __proto__: [] } instanceof Array &&
                      function (e, t) {
                        e.__proto__ = t
                      }) ||
                    function (e, t) {
                      for (var r in t)
                        Object.prototype.hasOwnProperty.call(t, r) &&
                          (e[r] = t[r])
                    }),
                  R(e, t)
                )
              }),
              function (e, t) {
                if ("function" != typeof t && null !== t)
                  throw new TypeError(
                    "Class extends value " +
                      String(t) +
                      " is not a constructor or null",
                  )
                function r() {
                  this.constructor = e
                }
                R(e, t),
                  (e.prototype =
                    null === t
                      ? Object.create(t)
                      : ((r.prototype = t.prototype), new r()))
              }),
            C = function () {
              return (
                (C =
                  Object.assign ||
                  function (e) {
                    for (var t, r = 1, n = arguments.length; r < n; r++)
                      for (var o in (t = arguments[r]))
                        Object.prototype.hasOwnProperty.call(t, o) &&
                          (e[o] = t[o])
                    return e
                  }),
                C.apply(this, arguments)
              )
            },
            j = (function (e) {
              function t() {
                return (null !== e && e.apply(this, arguments)) || this
              }
              return (
                x(t, e),
                (t.prototype.handleErrors = function (e) {
                  var t = this,
                    r = window.onerror
                  window.onerror = function (n, a, c, s, u) {
                    r && r(n, a, c, s, u), (n = I(e, n))
                    var l = {
                      category: o.JS_ERROR,
                      grade: i.ERROR,
                      errorUrl: a,
                      line: c,
                      col: s,
                      message: n,
                      stack: (u && u.stack) || "",
                    }
                    if (P(e, l)) {
                      var f = T(e, [])
                      ;(t.logInfo = C(C({}, f), l)), t.traceInfo(e.sendErrorTime)
                    }
                  }
                }),
                t
              )
            })(O),
            H = new j(),
            M = (function () {
              var e = function (t, r) {
                return (
                  (e =
                    Object.setPrototypeOf ||
                    ({ __proto__: [] } instanceof Array &&
                      function (e, t) {
                        e.__proto__ = t
                      }) ||
                    function (e, t) {
                      for (var r in t)
                        Object.prototype.hasOwnProperty.call(t, r) &&
                          (e[r] = t[r])
                    }),
                  e(t, r)
                )
              }
              return function (t, r) {
                if ("function" != typeof r && null !== r)
                  throw new TypeError(
                    "Class extends value " +
                      String(r) +
                      " is not a constructor or null",
                  )
                function n() {
                  this.constructor = t
                }
                e(t, r),
                  (t.prototype =
                    null === r
                      ? Object.create(r)
                      : ((n.prototype = r.prototype), new n()))
              }
            })(),
            D = function () {
              return (
                (D =
                  Object.assign ||
                  function (e) {
                    for (var t, r = 1, n = arguments.length; r < n; r++)
                      for (var o in (t = arguments[r]))
                        Object.prototype.hasOwnProperty.call(t, o) &&
                          (e[o] = t[o])
                    return e
                  }),
                D.apply(this, arguments)
              )
            },
            A = (function (e) {
              function t() {
                return (null !== e && e.apply(this, arguments)) || this
              }
              return (
                M(t, e),
                (t.prototype.handleErrors = function (e) {
                  var t = this
                  window.addEventListener("unhandledrejection", function (r) {
                    var n
                    if (
                      r &&
                      r.reason &&
                      !(null === (n = null == r ? void 0 : r.reason) ||
                      void 0 === n
                        ? void 0
                        : n.request)
                    ) {
                      var a = ""
                      try {
                        r.reason.config &&
                          r.reason.config.url &&
                          (a = r.reason.config.url)
                        var c = r.reason.message || ""
                        c = I(e, c)
                        var s = {
                          category: o.PROMISE_ERROR,
                          grade: i.ERROR,
                          errorUrl: a || location.href,
                          message: c,
                          stack: r.reason.stack,
                        }
                        if (!P(e, s)) return
                        var u = T(e, [])
                        ;(t.logInfo = D(D({}, u), s)),
                          t.traceInfo(e.sendErrorTime)
                      } catch (e) {
                        console.log(e)
                      }
                    }
                  })
                }),
                t
              )
            })(O),
            q = new A(),
            L = (function () {
              var e = function (t, r) {
                return (
                  (e =
                    Object.setPrototypeOf ||
                    ({ __proto__: [] } instanceof Array &&
                      function (e, t) {
                        e.__proto__ = t
                      }) ||
                    function (e, t) {
                      for (var r in t)
                        Object.prototype.hasOwnProperty.call(t, r) &&
                          (e[r] = t[r])
                    }),
                  e(t, r)
                )
              }
              return function (t, r) {
                if ("function" != typeof r && null !== r)
                  throw new TypeError(
                    "Class extends value " +
                      String(r) +
                      " is not a constructor or null",
                  )
                function n() {
                  this.constructor = t
                }
                e(t, r),
                  (t.prototype =
                    null === r
                      ? Object.create(r)
                      : ((n.prototype = r.prototype), new n()))
              }
            })(),
            N = function () {
              return (
                (N =
                  Object.assign ||
                  function (e) {
                    for (var t, r = 1, n = arguments.length; r < n; r++)
                      for (var o in (t = arguments[r]))
                        Object.prototype.hasOwnProperty.call(t, o) &&
                          (e[o] = t[o])
                    return e
                  }),
                N.apply(this, arguments)
              )
            },
            V = (function (e) {
              function t() {
                return (null !== e && e.apply(this, arguments)) || this
              }
              return (
                L(t, e),
                (t.prototype.handleError = function (e) {
                  var t = this
                  if (window.XMLHttpRequest) {
                    var r = XMLHttpRequest.prototype.send,
                      n = function (r) {
                        var n = r.start
                        return function (r) {
                          var c,
                            s,
                            u,
                            l,
                            f,
                            p,
                            d = Date.now() - n,
                            h = d >= e.maxReqCostTime,
                            v = r.type
                          try {
                            if (
                              !h &&
                              !(
                                (null ===
                                  (c = null == r ? void 0 : r.currentTarget) ||
                                void 0 === c
                                  ? void 0
                                  : c.status) >= 400 ||
                                0 ===
                                  (null ===
                                    (s = null == r ? void 0 : r.currentTarget) ||
                                  void 0 === s
                                    ? void 0
                                    : s.status)
                              )
                            )
                              return
                            var g = r.currentTarget.status,
                              m = void 0
                            "timeout" === v
                              ? ((m = "net::timeout (status: ".concat(g, ")")),
                                (null === (u = null == r ? void 0 : r.target) ||
                                void 0 === u
                                  ? void 0
                                  : u.timeout) &&
                                  (m += ", ajax timeout: ".concat(
                                    r.target.timeout,
                                  )))
                              : (m = h
                                  ? "net::slowRequest (status: "
                                      .concat(g, "), maxReqCostTime: ")
                                      .concat(e.maxReqCostTime)
                                  : "net::ERR_EMPTY(status: ".concat(g, ")"))
                            var y = e.collector + a.ERRORS
                            if (
                              (null === (l = null == r ? void 0 : r.target) ||
                              void 0 === l
                                ? void 0
                                : l.getRequestConfig) &&
                              r.target.getRequestConfig[1] === y
                            )
                              return
                            var w = r.target.response
                            ;(null == w ? void 0 : w.length) > 150 &&
                              (w = w.replace(/\s/g, "")).length > 150 &&
                              (w = ""
                                .concat(w.slice(0, 150), "...(")
                                .concat(g, ")"))
                            var E = "接口响应正常"
                            "load" !== v && (E = I(e, (E = w || m)))
                            var O = {
                              category: o.AJAX_ERROR,
                              grade: i.ERROR,
                              errorUrl: (
                                null === (f = r.target) || void 0 === f
                                  ? void 0
                                  : f.getRequestConfig
                              )
                                ? null === (p = r.target) || void 0 === p
                                  ? void 0
                                  : p.getRequestConfig[1]
                                : "",
                              message: E,
                              reqErrType: v,
                              stack: m,
                            }
                            if (
                              ((h || "timeout" === v) && (O.slowRequest = d),
                              "load" !== v &&
                                (O.stack = r.type + ": " + (w || m)),
                              !P(e, O))
                            )
                              return
                            var b = T(e, [])
                            t.logInfo = N(N({}, b), O)
                            var R = sessionStorage.getItem("current-sw-traceId")
                            sessionStorage.removeItem("current-sw-traceId"),
                              R && (t.logInfo.traceId = R),
                              t.traceInfo(e.sendErrorTime)
                          } catch (e) {
                            console.log(e)
                          }
                        }
                      }
                    XMLHttpRequest.prototype.send = function () {
                      var e = Date.now()
                      if (this.addEventListener)
                        this.addEventListener("load", n({ start: e })),
                          this.addEventListener("error", n({ start: e })),
                          this.addEventListener("abort", n({ start: e })),
                          this.addEventListener("timeout", n({ start: e }))
                      else {
                        var t = this.onreadystatechange
                        this.onreadystatechange = function (r) {
                          t.apply(this, arguments),
                            4 === this.readyState && n({ start: e })(r)
                        }
                      }
                      return r.apply(this, arguments)
                    }
                  }
                }),
                t
              )
            })(O),
            G = new V(),
            U = (function () {
              var e = function (t, r) {
                return (
                  (e =
                    Object.setPrototypeOf ||
                    ({ __proto__: [] } instanceof Array &&
                      function (e, t) {
                        e.__proto__ = t
                      }) ||
                    function (e, t) {
                      for (var r in t)
                        Object.prototype.hasOwnProperty.call(t, r) &&
                          (e[r] = t[r])
                    }),
                  e(t, r)
                )
              }
              return function (t, r) {
                if ("function" != typeof r && null !== r)
                  throw new TypeError(
                    "Class extends value " +
                      String(r) +
                      " is not a constructor or null",
                  )
                function n() {
                  this.constructor = t
                }
                e(t, r),
                  (t.prototype =
                    null === r
                      ? Object.create(r)
                      : ((n.prototype = r.prototype), new n()))
              }
            })(),
            F = function () {
              return (
                (F =
                  Object.assign ||
                  function (e) {
                    for (var t, r = 1, n = arguments.length; r < n; r++)
                      for (var o in (t = arguments[r]))
                        Object.prototype.hasOwnProperty.call(t, o) &&
                          (e[o] = t[o])
                    return e
                  }),
                F.apply(this, arguments)
              )
            },
            W = function (e, t, r) {
              if (r || 2 === arguments.length)
                for (var n, o = 0, i = t.length; o < i; o++)
                  (!n && o in t) ||
                    (n || (n = Array.prototype.slice.call(t, 0, o)),
                    (n[o] = t[o]))
              return e.concat(n || Array.prototype.slice.call(t))
            },
            B = (function (e) {
              function t() {
                return (null !== e && e.apply(this, arguments)) || this
              }
              return (
                U(t, e),
                (t.prototype.handleErrors = function (e) {
                  var t = this
                  window.addEventListener(
                    "error",
                    function (r) {
                      if (r)
                        try {
                          var n = r.target || r.srcElement
                          if (
                            !(
                              n instanceof HTMLScriptElement ||
                              n instanceof HTMLLinkElement ||
                              n instanceof HTMLImageElement
                            )
                          )
                            return
                          var a = n.src || n.href || ""
                          if (!a || a === location.href) return
                          if (
                            W(
                              W(
                                [],
                                [
                                  "stat.10jqka.com.cn",
                                  "jjstat.10jqka.com.cn",
                                  "razor.10jqka.com.cn",
                                ],
                                !0,
                              ),
                              e.noErrorPaths,
                              !0,
                            ).some(function (e) {
                              return a.indexOf(e) > -1
                            })
                          )
                            return
                          var c = "load ".concat(n.tagName, " resource error")
                          c = I(e, c)
                          var s = {
                            category: o.RESOURCE_ERROR,
                            grade: "IMG" === n.tagName ? i.WARNING : i.ERROR,
                            errorUrl: a,
                            message: c,
                            stack: "load ".concat(
                              n.src || n.href,
                              " resource error",
                            ),
                          }
                          if (!P(e, s)) return
                          var u = T(e, [])
                          ;(t.logInfo = F(F({}, u), s)),
                            t.traceInfo(e.sendErrorTime)
                        } catch (e) {
                          throw e
                        }
                    },
                    !0,
                  )
                }),
                t
              )
            })(O),
            X = new B(),
            J = (function () {
              var e = function (t, r) {
                return (
                  (e =
                    Object.setPrototypeOf ||
                    ({ __proto__: [] } instanceof Array &&
                      function (e, t) {
                        e.__proto__ = t
                      }) ||
                    function (e, t) {
                      for (var r in t)
                        Object.prototype.hasOwnProperty.call(t, r) &&
                          (e[r] = t[r])
                    }),
                  e(t, r)
                )
              }
              return function (t, r) {
                if ("function" != typeof r && null !== r)
                  throw new TypeError(
                    "Class extends value " +
                      String(r) +
                      " is not a constructor or null",
                  )
                function n() {
                  this.constructor = t
                }
                e(t, r),
                  (t.prototype =
                    null === r
                      ? Object.create(r)
                      : ((n.prototype = r.prototype), new n()))
              }
            })(),
            K = function () {
              return (
                (K =
                  Object.assign ||
                  function (e) {
                    for (var t, r = 1, n = arguments.length; r < n; r++)
                      for (var o in (t = arguments[r]))
                        Object.prototype.hasOwnProperty.call(t, o) &&
                          (e[o] = t[o])
                    return e
                  }),
                K.apply(this, arguments)
              )
            },
            Q = (function (e) {
              function t() {
                return (null !== e && e.apply(this, arguments)) || this
              }
              return (
                J(t, e),
                (t.prototype.handleErrors = function (e, t) {
                  var r = this
                  if (!t)
                    return console.warn(
                      "Vue handle errors set failed.(SkyWalking)",
                    )
                  t.config.errorHandler = function (t, n, a) {
                    try {
                      var c = a
                      c = I(e, c)
                      var s = {
                        category: o.VUE_ERROR,
                        grade: i.ERROR,
                        errorUrl: location.href,
                        message: c,
                        stack: t.stack,
                      }
                      if (!P(e, s)) return
                      var u = T(e, [])
                      ;(r.logInfo = K(K({}, u), s)), r.traceInfo(e.sendErrorTime)
                    } catch (e) {
                      throw e
                    }
                  }
                }),
                t
              )
            })(O),
            z = new Q(),
            Y = (function () {
              var e = function (t, r) {
                return (
                  (e =
                    Object.setPrototypeOf ||
                    ({ __proto__: [] } instanceof Array &&
                      function (e, t) {
                        e.__proto__ = t
                      }) ||
                    function (e, t) {
                      for (var r in t)
                        Object.prototype.hasOwnProperty.call(t, r) &&
                          (e[r] = t[r])
                    }),
                  e(t, r)
                )
              }
              return function (t, r) {
                if ("function" != typeof r && null !== r)
                  throw new TypeError(
                    "Class extends value " +
                      String(r) +
                      " is not a constructor or null",
                  )
                function n() {
                  this.constructor = t
                }
                e(t, r),
                  (t.prototype =
                    null === r
                      ? Object.create(r)
                      : ((n.prototype = r.prototype), new n()))
              }
            })(),
            $ = function () {
              return (
                ($ =
                  Object.assign ||
                  function (e) {
                    for (var t, r = 1, n = arguments.length; r < n; r++)
                      for (var o in (t = arguments[r]))
                        Object.prototype.hasOwnProperty.call(t, o) &&
                          (e[o] = t[o])
                    return e
                  }),
                $.apply(this, arguments)
              )
            },
            Z = (function (e) {
              function t() {
                return (null !== e && e.apply(this, arguments)) || this
              }
              return (
                Y(t, e),
                (t.prototype.handleErrors = function (e, t) {
                  t = t || { name: "", message: "", stack: "" }
                  var r = e.category,
                    n = e.grade,
                    a = e.level,
                    c = e.content,
                    s = e.timeout,
                    u = {
                      category: r || o.JS_ERROR,
                      grade: n || i.ERROR,
                      errorUrl: t.name || location.href,
                      message: t.message,
                      stack: t.stack,
                    }
                  if (
                    ("number" == typeof a && (u.level = a),
                    c && (u.content = c),
                    P(e, u))
                  ) {
                    var l = T(e, [])
                    ;(this.logInfo = $($({}, l), u)), this.traceInfo(s)
                  }
                }),
                t
              )
            })(O),
            ee = new Z(),
            te = (function () {
              function e() {}
              return (
                (e.prototype.getPerfTiming = function (e) {
                  try {
                    return this.returnPerfData(e)
                  } catch (e) {
                    throw e
                  }
                }),
                (e.prototype.returnPerfData = function (e) {
                  if (window.performance && window.performance.timing) {
                    var r = window.performance.timing,
                      n = 0
                    void 0 !== r.navigationStart
                      ? (n = t(r.fetchStart - r.navigationStart))
                      : void 0 !== r.redirectEnd &&
                        (n = t(r.redirectEnd - r.redirectStart))
                    var o = {
                      redirectTime: n,
                      dnsTime: t(r.domainLookupEnd - r.domainLookupStart),
                      tcpTime: t(r.connectEnd - r.connectStart),
                      ttfbTime: t(r.responseStart - r.requestStart),
                      transTime: t(r.responseEnd - r.responseStart),
                      domAnalysisTime: t(r.domInteractive - r.responseEnd),
                      fptTime: t(r.responseEnd - r.fetchStart),
                      domReadyTime: t(r.domContentLoadedEventEnd - r.fetchStart),
                      resTime: t(r.loadEventStart - r.domContentLoadedEventEnd),
                      loadPageTime: t(r.loadEventStart - r.fetchStart),
                      sslTime:
                        "https:" === location.protocol &&
                        r.secureConnectionStart > 0
                          ? t(r.connectEnd - r.secureConnectionStart)
                          : void 0,
                      ttlTime: t(r.domInteractive - r.fetchStart),
                      firstPackTime: t(r.responseStart - r.domainLookupStart),
                      fmpTime: 0,
                      sdkVersion: k,
                    }
                    if (e) {
                      var i = ["script", "link", "img"],
                        a = []
                      performance.getEntries().forEach(function (e, t) {
                        if (t) {
                          var r = e.name,
                            n = e.initiatorType,
                            o = e.transferSize,
                            c = e.responseEnd,
                            s = e.startTime
                          if (-1 !== i.indexOf(n)) {
                            var u = Math.round(c - s),
                              l = r.split("/").pop().replace(/^;/, "")
                            a.push({ name: l, type: n, size: o, time: u })
                          }
                        }
                      }),
                        (o.entries = a)
                    }
                    return o
                  }
                  console.log("your browser do not support performance")
                }),
                e
              )
            })(),
            re = function (e, t) {
              return window.getComputedStyle
                ? window.getComputedStyle(e, null)[t]
                : e.currentStyle[t]
            }
          !(function (e) {
            ;(e[(e.SVG = 2)] = "SVG"),
              (e[(e.IMG = 2)] = "IMG"),
              (e[(e.CANVAS = 4)] = "CANVAS"),
              (e[(e.OBJECT = 4)] = "OBJECT"),
              (e[(e.EMBED = 4)] = "EMBED"),
              (e[(e.VIDEO = 4)] = "VIDEO")
          })(S || (S = {}))
          var ne = performance.now(),
            oe = ["SCRIPT", "STYLE", "META", "HEAD", "LINK"],
            ie = window.innerWidth,
            ae = window.innerHeight,
            ce = (function () {
              function e() {
                ;(this.fmpTime = 0),
                  (this.statusCollector = []),
                  (this.flag = !0),
                  (this.observer = null),
                  (this.callbackCount = 0),
                  (this.entries = {}),
                  performance && performance.getEntries
                    ? this.initObserver()
                    : console.log(
                        "your browser do not support performance.getEntries",
                      )
              }
              return (
                (e.prototype.getFirstSnapShot = function () {
                  var e = performance.now(),
                    t = document.body
                  t && this.setTag(t, this.callbackCount),
                    this.statusCollector.push({ time: e })
                }),
                (e.prototype.initObserver = function () {
                  var e = this
                  this.getFirstSnapShot(),
                    (this.observer = new MutationObserver(function () {
                      e.callbackCount += 1
                      var t = performance.now(),
                        r = document.body
                      r && e.setTag(r, e.callbackCount),
                        e.statusCollector.push({ time: t })
                    })),
                    this.observer.observe(document, {
                      childList: !0,
                      subtree: !0,
                    }),
                    "complete" === document.readyState
                      ? this.calculateFinalScore()
                      : window.addEventListener(
                          "load",
                          function () {
                            e.calculateFinalScore()
                          },
                          !1,
                        )
                }),
                (e.prototype.calculateFinalScore = function () {
                  var e = this
                  if (MutationEvent && this.flag)
                    if (this.checkNeedCancel(ne)) {
                      this.observer.disconnect(), (this.flag = !1)
                      for (
                        var t = null,
                          r = 0,
                          n = this.getTreeScore(document.body).dpss;
                        r < n.length;
                        r++
                      ) {
                        var o = n[r]
                        t && t.st ? t.st < o.st && (t = o) : (t = o)
                      }
                      if (
                        (performance.getEntries().forEach(function (t) {
                          e.entries[t.name] = t.responseEnd
                        }),
                        !t)
                      )
                        return !1
                      var i = this.filterResult(t.els),
                        a = this.getFmpTime(i)
                      this.fmpTime = a
                    } else
                      setTimeout(function () {
                        e.calculateFinalScore()
                      }, 500)
                }),
                (e.prototype.getFmpTime = function (e) {
                  var t = 0
                  try {
                    for (var r = 0, n = e; r < n.length; r++) {
                      var o = n[r],
                        i = 0
                      if (1 === o.weight) {
                        var a = parseInt(o.ele.getAttribute("fmp_c"), 10)
                        i = this.statusCollector[a].time
                      } else if (2 === o.weight)
                        if ("IMG" === o.ele.tagName) i = this.entries[o.ele.src]
                        else if ("SVG" === o.ele.tagName) {
                          a = parseInt(o.ele.getAttribute("fmp_c"), 10)
                          i = this.statusCollector[a].time
                        } else {
                          var c = re(o.ele, "background-image").match(
                              /url\(\"(.*?)\"\)/,
                            ),
                            s = void 0
                          c && c[1] && (s = c[1]),
                            (null == s ? void 0 : s.indexOf("http")) > -1 &&
                              (s = location.protocol + ((c && c[1]) || "-")),
                            (i = this.entries[s])
                        }
                      else if (4 === o.weight)
                        if ("CANVAS" === o.ele.tagName) {
                          a = parseInt(o.ele.getAttribute("fmp_c"), 10)
                          i =
                            this.statusCollector[a] &&
                            this.statusCollector[a].time
                        } else
                          "VIDEO" === o.ele.tagName &&
                            ((i = this.entries[o.ele.src]) ||
                              (i = this.entries[o.ele.poster]))
                      "number" != typeof i && (i = 0), t < i && (t = i)
                    }
                  } catch (e) {
                    console.warn(e)
                  }
                  return t
                }),
                (e.prototype.filterResult = function (e) {
                  if (1 === e.length) return e
                  var t = 0
                  e.forEach(function (e) {
                    t += e.st
                  })
                  var r = t / e.length
                  return e.filter(function (e) {
                    return e.st > r
                  })
                }),
                (e.prototype.checkNeedCancel = function (e) {
                  var t = performance.now() - e,
                    r =
                      this.statusCollector.length > 0
                        ? this.statusCollector[this.statusCollector.length - 1]
                            .time
                        : 0
                  return t > 3e3 || t - r > 1e3
                }),
                (e.prototype.getTreeScore = function (e) {
                  if (!e) return {}
                  for (var t = [], r = 0, n = e.children; r < n.length; r++) {
                    var o = n[r]
                    if (o.getAttribute("fmp_c")) {
                      var i = this.getTreeScore(o)
                      i.st && t.push(i)
                    }
                  }
                  return this.calcaulteGrades(e, t)
                }),
                (e.prototype.calcaulteGrades = function (e, t) {
                  var r = e.getBoundingClientRect(),
                    n = r.width,
                    o = r.height,
                    i = r.left,
                    a = r.top,
                    c = !0
                  ;(ae < a || ie < i) && (c = !1)
                  var s = 0
                  t.forEach(function (e) {
                    s += e.st
                  })
                  var u = Number(S[e.tagName]) || 1
                  1 === u &&
                    re(e, "background-image") &&
                    "initial" !== re(e, "background-image") &&
                    "none" !== re(e, "background-image") &&
                    (u = S.IMG)
                  var l = c ? n * o * u : 0,
                    f = [{ ele: e, st: l, weight: u }],
                    p = e,
                    d = this.calculateAreaParent(e)
                  if (s > l * d || 0 === d) {
                    ;(l = s), (f = [])
                    for (var h = 0, v = t; h < v.length; h++) {
                      var g = v[h]
                      f = f.concat(g.els)
                    }
                  }
                  return { dpss: t, st: l, els: f, root: p }
                }),
                (e.prototype.calculateAreaParent = function (e) {
                  var t = e.getBoundingClientRect(),
                    r = t.left,
                    n = t.right,
                    o = t.top,
                    i = t.bottom,
                    a = t.width,
                    c = t.height,
                    s = ie,
                    u = ae,
                    l = n - r + (s - 0) - (Math.max(n, s) - Math.min(r, 0)),
                    f = i - o + (u - 0) - (Math.max(i, u) - Math.min(o, 0))
                  return l <= 0 || f <= 0 ? 0 : (l * f) / (a * c)
                }),
                (e.prototype.setTag = function (e, t) {
                  var r = e.tagName
                  if (-1 === oe.indexOf(r)) {
                    var n = e.children
                    if (n && n.length > 0)
                      for (var o = n.length - 1; o >= 0; o--) {
                        var i = n[o]
                        if (!(null !== i.getAttribute("fmp_c"))) {
                          var a = i.getBoundingClientRect(),
                            c = a.left,
                            s = a.top,
                            u = a.width,
                            l = a.height
                          if (ae < s || ie < c || 0 === u || 0 === l) continue
                          i.setAttribute("fmp_c", "".concat(t))
                        }
                        this.setTag(i, t)
                      }
                  }
                }),
                e
              )
            })(),
            se = function () {
              return (
                (se =
                  Object.assign ||
                  function (e) {
                    for (var t, r = 1, n = arguments.length; r < n; r++)
                      for (var o in (t = arguments[r]))
                        Object.prototype.hasOwnProperty.call(t, o) &&
                          (e[o] = t[o])
                    return e
                  }),
                se.apply(this, arguments)
              )
            },
            ue = function (e, t, r, n) {
              return new (r || (r = Promise))(function (o, i) {
                function a(e) {
                  try {
                    s(n.next(e))
                  } catch (e) {
                    i(e)
                  }
                }
                function c(e) {
                  try {
                    s(n.throw(e))
                  } catch (e) {
                    i(e)
                  }
                }
                function s(e) {
                  var t
                  e.done
                    ? o(e.value)
                    : ((t = e.value),
                      t instanceof r
                        ? t
                        : new r(function (e) {
                            e(t)
                          })).then(a, c)
                }
                s((n = n.apply(e, t || [])).next())
              })
            },
            le = function (e, t) {
              var r,
                n,
                o,
                i,
                a = {
                  label: 0,
                  sent: function () {
                    if (1 & o[0]) throw o[1]
                    return o[1]
                  },
                  trys: [],
                  ops: [],
                }
              return (
                (i = { next: c(0), throw: c(1), return: c(2) }),
                "function" == typeof Symbol &&
                  (i[Symbol.iterator] = function () {
                    return this
                  }),
                i
              )
              function c(c) {
                return function (s) {
                  return (function (c) {
                    if (r) throw new TypeError("Generator is already executing.")
                    for (; i && ((i = 0), c[0] && (a = 0)), a; )
                      try {
                        if (
                          ((r = 1),
                          n &&
                            (o =
                              2 & c[0]
                                ? n.return
                                : c[0]
                                  ? n.throw || ((o = n.return) && o.call(n), 0)
                                  : n.next) &&
                            !(o = o.call(n, c[1])).done)
                        )
                          return o
                        switch (((n = 0), o && (c = [2 & c[0], o.value]), c[0])) {
                          case 0:
                          case 1:
                            o = c
                            break
                          case 4:
                            return a.label++, { value: c[1], done: !1 }
                          case 5:
                            a.label++, (n = c[1]), (c = [0])
                            continue
                          case 7:
                            ;(c = a.ops.pop()), a.trys.pop()
                            continue
                          default:
                            if (
                              !((o = a.trys),
                              (o = o.length > 0 && o[o.length - 1]) ||
                                (6 !== c[0] && 2 !== c[0]))
                            ) {
                              a = 0
                              continue
                            }
                            if (
                              3 === c[0] &&
                              (!o || (c[1] > o[0] && c[1] < o[3]))
                            ) {
                              a.label = c[1]
                              break
                            }
                            if (6 === c[0] && a.label < o[1]) {
                              ;(a.label = o[1]), (o = c)
                              break
                            }
                            if (o && a.label < o[2]) {
                              ;(a.label = o[2]), a.ops.push(c)
                              break
                            }
                            o[2] && a.ops.pop(), a.trys.pop()
                            continue
                        }
                        c = t.call(e, a)
                      } catch (e) {
                        ;(c = [6, e]), (n = 0)
                      } finally {
                        r = o = 0
                      }
                    if (5 & c[0]) throw c[1]
                    return { value: c[0] ? c[1] : void 0, done: !0 }
                  })([c, s])
                }
              }
            },
            fe = new ((function () {
              function e() {
                this.perfConfig = { perfDetail: {} }
              }
              return (
                (e.prototype.recordPerf = function (e, t) {
                  return ue(this, void 0, void 0, function () {
                    var r,
                      o,
                      i = this
                    return le(this, function (a) {
                      switch (a.label) {
                        case 0:
                          return (
                            (r = { fmpTime: void 0 }),
                            e.autoTracePerf
                              ? ((o = this.perfConfig),
                                [4, new te().getPerfTiming(e.perfEntries)])
                              : [3, 3]
                          )
                        case 1:
                          return (
                            (o.perfDetail = a.sent()),
                            e.useFmp ? [4, new ce()] : [3, 3]
                          )
                        case 2:
                          ;(r = a.sent()), (a.label = 3)
                        case 3:
                          return (
                            setTimeout(function () {
                              var o,
                                a = e.autoTracePerf
                                  ? se(se({}, i.perfConfig.perfDetail), {
                                      fmpTime: e.useFmp
                                        ? parseInt(String(r.fmpTime), 10)
                                        : void 0,
                                    })
                                  : void 0,
                                c = e.pagePath,
                                s = e.group,
                                u = e.subGroup,
                                l = e.service,
                                f = e.serviceVersion,
                                p = c
                              p ||
                                (p =
                                  location.href.replace(location.origin, "") ||
                                  "/")
                              var d = n(),
                                h = d.year,
                                v = d.month,
                                g = d.day,
                                m = d.hours,
                                w = d.minutes,
                                E = d.seconds,
                                O = d.miSeconds,
                                b = ""
                                  .concat(h, "-")
                                  .concat(v, "-")
                                  .concat(g, " ")
                                  .concat(m, ":")
                                  .concat(w, ":")
                                  .concat(E, ".")
                                  .concat(O),
                                R = se(se({}, a), {
                                  pagePath: p,
                                  group: s,
                                  subGroup: u,
                                  service: l,
                                  serviceVersion: f,
                                  logtime: b,
                                })
                              t && (R.uniqueId = t)
                              var S =
                                (null === (o = e.hooks) || void 0 === o
                                  ? void 0
                                  : o.reportHook) || e.reportHook
                              new y("PERF", e.collector, S).sendByXhr(R),
                                i.clearPerf()
                            }, e.sendPerfTime || 6e3),
                            [2]
                          )
                      }
                    })
                  })
                }),
                (e.prototype.clearPerf = function () {
                  window.performance &&
                    window.performance.clearResourceTimings &&
                    (window.performance.clearResourceTimings(),
                    (this.perfConfig = { perfDetail: {} }))
                }),
                e
              )
            })())(),
            pe = (function () {
              var e = function (t, r) {
                return (
                  (e =
                    Object.setPrototypeOf ||
                    ({ __proto__: [] } instanceof Array &&
                      function (e, t) {
                        e.__proto__ = t
                      }) ||
                    function (e, t) {
                      for (var r in t)
                        Object.prototype.hasOwnProperty.call(t, r) &&
                          (e[r] = t[r])
                    }),
                  e(t, r)
                )
              }
              return function (t, r) {
                if ("function" != typeof r && null !== r)
                  throw new TypeError(
                    "Class extends value " +
                      String(r) +
                      " is not a constructor or null",
                  )
                function n() {
                  this.constructor = t
                }
                e(t, r),
                  (t.prototype =
                    null === r
                      ? Object.create(r)
                      : ((n.prototype = r.prototype), new n()))
              }
            })(),
            de = function () {
              return (
                (de =
                  Object.assign ||
                  function (e) {
                    for (var t, r = 1, n = arguments.length; r < n; r++)
                      for (var o in (t = arguments[r]))
                        Object.prototype.hasOwnProperty.call(t, o) &&
                          (e[o] = t[o])
                    return e
                  }),
                de.apply(this, arguments)
              )
            },
            he = (function (e) {
              function t() {
                var t = (null !== e && e.apply(this, arguments)) || this
                return (t.isLockMonitor = !1), t
              }
              return (
                pe(t, e),
                (t.prototype.report = function (e, t) {
                  var r = {
                      message: "检测出白屏，检测方式为: ".concat(t),
                      grade: i.ERROR,
                      category: o.WHITE_SCREEN_ERROR,
                      errorUrl: location.href,
                      whiteScreen: t,
                    },
                    n = T(e, [])
                  this.logInfo = de(de({}, n), r)
                  this.traceInfo(10),
                    "timeout" === t &&
                      fe.recordPerf(
                        de(de(de({}, e), this.logInfo), { sendPerfTime: 10 }),
                        n.uniqueId,
                      )
                }),
                (t.prototype.monitorWhiteScreen = function (e, t) {
                  var r, n
                  if (!this.isLockMonitor) {
                    this.isLockMonitor = !0
                    var o = !1,
                      i =
                        document.querySelector("#app") ||
                        document.querySelector("#root")
                    if (i)
                      (i.innerHTML
                        ? i.innerHTML.replace(/<!---->/g, "").trim()
                        : ""
                      ).length || (o = !0)
                    else if (
                      !(null === (r = document.querySelector("body")) ||
                      void 0 === r
                        ? void 0
                        : r.innerText)
                    ) {
                      o = !0
                      for (
                        var a = document.querySelectorAll("body > div"),
                          c = a.length,
                          s = 0;
                        s < c;
                        s++
                      )
                        if (
                          null === (n = a[s]) || void 0 === n
                            ? void 0
                            : n.innerHTML
                        ) {
                          o = !1
                          break
                        }
                    }
                    o &&
                      (console.log("".concat(t, " 检测出白屏了")),
                      this.report(e, t))
                  }
                }),
                (t.prototype.handleWhiteScreen = function (e) {
                  var t = this
                  setTimeout(function () {
                    t.monitorWhiteScreen(e, "timeout")
                  }, 2500),
                    window.addEventListener("load", function () {
                      setTimeout(function () {
                        t.monitorWhiteScreen(e, "load")
                      }, 10)
                    })
                }),
                t
              )
            })(O),
            ve = new he(),
            ge = function () {
              return (
                (ge =
                  Object.assign ||
                  function (e) {
                    for (var t, r = 1, n = arguments.length; r < n; r++)
                      for (var o in (t = arguments[r]))
                        Object.prototype.hasOwnProperty.call(t, o) &&
                          (e[o] = t[o])
                    return e
                  }),
                ge.apply(this, arguments)
              )
            },
            me = function (e) {
              return window.btoa(e)
            }
          function ye(e, t) {
            var r = window.XMLHttpRequest,
              n = XMLHttpRequest.prototype.send,
              o = XMLHttpRequest.prototype.open
            function i(e) {
              var t = new CustomEvent(e, { detail: this })
              window.dispatchEvent(t)
            }
            ;(r.getRequestConfig = []),
              (window.XMLHttpRequest = function () {
                var e = new r()
                return (
                  e.addEventListener(
                    "readystatechange",
                    function () {
                      i.call(this, "xhrReadyStateChange")
                    },
                    !1,
                  ),
                  (e.open = function (e, t, r, n, i) {
                    return (
                      (this.getRequestConfig = arguments),
                      o.apply(this, arguments)
                    )
                  }),
                  (e.send = function (e) {
                    return n.apply(this, arguments)
                  }),
                  e
                )
              })
            var c = []
            window.addEventListener("xhrReadyStateChange", function (r) {
              var n,
                o = {
                  traceId: "",
                  service: e.service + d,
                  spans: [],
                  serviceInstance: e.serviceVersion,
                  traceSegmentId: "",
                },
                i = r.detail.readyState,
                s = r.detail.getRequestConfig,
                h = s[1],
                v = {}
              0 === s[1].indexOf("http://") || 0 === s[1].indexOf("https://")
                ? (v = new URL(s[1]))
                : 0 === s[1].indexOf("//")
                  ? (v = new URL(
                      "".concat(window.location.protocol).concat(s[1]),
                    ))
                  : ((v = new URL(window.location.href)).pathname = s[1])
              if (
                (!(null === (n = e.noTraceUrls) || void 0 === n
                  ? void 0
                  : n.length) ||
                  !e.noTraceUrls.some(function (e) {
                    if ("string" == typeof e) {
                      if (e === h) return !0
                    } else if (e instanceof RegExp && e.test(h)) return !0
                  })) &&
                (!(
                  [a.ERROR, a.ERRORS, a.PERF, a.SEGMENTS].indexOf(v.pathname) > -1
                ) ||
                  e.traceSDKInternal)
              ) {
                if (i === u.OPENED) {
                  var g = b(),
                    m = b()
                  sessionStorage.setItem("current-sw-traceId", g),
                    c.push({
                      event: r.detail,
                      startTime: new Date().getTime(),
                      traceId: g,
                      traceSegmentId: m,
                    })
                  var y = window.ELK_WHITESCREEN_TRACEID || String(me(g)),
                    w = String(me(m)),
                    E = String(me(o.service)),
                    O = String(me(o.serviceInstance)),
                    R = String(me(e.pagePath)),
                    S = String(me(v.host)),
                    _ = o.spans.length,
                    k = ""
                      .concat(1, "-")
                      .concat(y, "-")
                      .concat(w, "-")
                      .concat(_, "-")
                      .concat(E, "-")
                      .concat(O, "-")
                      .concat(R, "-")
                      .concat(S)
                  r.detail.setRequestHeader("sw8", k)
                }
                if (i === u.DONE) {
                  for (var T = new Date().getTime(), I = 0; I < c.length; I++)
                    if (c[I].event.readyState === u.DONE) {
                      var P = {}
                      c[I].event.status && (P = new URL(c[I].event.responseURL))
                      var x = {
                        operationName: e.pagePath,
                        startTime: c[I].startTime,
                        endTime: T,
                        spanId: o.spans.length,
                        spanLayer: l,
                        spanType: f,
                        isError: 0 === r.detail.status || r.detail.status >= 400,
                        parentSpanId: o.spans.length - 1,
                        componentId: p,
                        peer: P.host,
                        tags: e.detailMode
                          ? [
                              { key: "http.method", value: s[0] },
                              { key: "url", value: c[I].event.responseURL },
                            ]
                          : void 0,
                      }
                      ;(o = ge(ge({}, o), {
                        traceId: c[I].traceId,
                        traceSegmentId: c[I].traceSegmentId,
                      })).spans.push(x),
                        c.splice(I, 1)
                    }
                  t.push(o)
                }
              }
            })
          }
          var we = function () {
              return (
                (we =
                  Object.assign ||
                  function (e) {
                    for (var t, r = 1, n = arguments.length; r < n; r++)
                      for (var o in (t = arguments[r]))
                        Object.prototype.hasOwnProperty.call(t, o) &&
                          (e[o] = t[o])
                    return e
                  }),
                we.apply(this, arguments)
              )
            },
            Ee = function (e, t, r, n) {
              return new (r || (r = Promise))(function (o, i) {
                function a(e) {
                  try {
                    s(n.next(e))
                  } catch (e) {
                    i(e)
                  }
                }
                function c(e) {
                  try {
                    s(n.throw(e))
                  } catch (e) {
                    i(e)
                  }
                }
                function s(e) {
                  var t
                  e.done
                    ? o(e.value)
                    : ((t = e.value),
                      t instanceof r
                        ? t
                        : new r(function (e) {
                            e(t)
                          })).then(a, c)
                }
                s((n = n.apply(e, t || [])).next())
              })
            },
            Oe = function (e, t) {
              var r,
                n,
                o,
                i,
                a = {
                  label: 0,
                  sent: function () {
                    if (1 & o[0]) throw o[1]
                    return o[1]
                  },
                  trys: [],
                  ops: [],
                }
              return (
                (i = { next: c(0), throw: c(1), return: c(2) }),
                "function" == typeof Symbol &&
                  (i[Symbol.iterator] = function () {
                    return this
                  }),
                i
              )
              function c(c) {
                return function (s) {
                  return (function (c) {
                    if (r) throw new TypeError("Generator is already executing.")
                    for (; i && ((i = 0), c[0] && (a = 0)), a; )
                      try {
                        if (
                          ((r = 1),
                          n &&
                            (o =
                              2 & c[0]
                                ? n.return
                                : c[0]
                                  ? n.throw || ((o = n.return) && o.call(n), 0)
                                  : n.next) &&
                            !(o = o.call(n, c[1])).done)
                        )
                          return o
                        switch (((n = 0), o && (c = [2 & c[0], o.value]), c[0])) {
                          case 0:
                          case 1:
                            o = c
                            break
                          case 4:
                            return a.label++, { value: c[1], done: !1 }
                          case 5:
                            a.label++, (n = c[1]), (c = [0])
                            continue
                          case 7:
                            ;(c = a.ops.pop()), a.trys.pop()
                            continue
                          default:
                            if (
                              !((o = a.trys),
                              (o = o.length > 0 && o[o.length - 1]) ||
                                (6 !== c[0] && 2 !== c[0]))
                            ) {
                              a = 0
                              continue
                            }
                            if (
                              3 === c[0] &&
                              (!o || (c[1] > o[0] && c[1] < o[3]))
                            ) {
                              a.label = c[1]
                              break
                            }
                            if (6 === c[0] && a.label < o[1]) {
                              ;(a.label = o[1]), (o = c)
                              break
                            }
                            if (o && a.label < o[2]) {
                              ;(a.label = o[2]), a.ops.push(c)
                              break
                            }
                            o[2] && a.ops.pop(), a.trys.pop()
                            continue
                        }
                        c = t.call(e, a)
                      } catch (e) {
                        ;(c = [6, e]), (n = 0)
                      } finally {
                        r = o = 0
                      }
                    if (5 & c[0]) throw c[1]
                    return { value: c[0] ? c[1] : void 0, done: !0 }
                  })([c, s])
                }
              }
            },
            be = function (e) {
              return window.btoa(e)
            }
          function Re(e) {
            var t,
              r = [],
              n =
                (null === (t = e.hooks) || void 0 === t
                  ? void 0
                  : t.reportHook) || e.reportHook
            e.noTrace ||
              (ye(e, r),
              (function (e, t) {
                var r = this,
                  n = window.fetch,
                  o = {
                    traceId: "",
                    service: e.service + d,
                    spans: [],
                    serviceInstance: e.serviceVersion,
                    traceSegmentId: "",
                  },
                  i = {}
                window.fetch = function () {
                  for (var c = [], s = 0; s < arguments.length; s++)
                    c[s] = arguments[s]
                  return (function (c) {
                    return Ee(r, void 0, void 0, function () {
                      var r, s, u, d, h, v, g, m, y, w, E, O, R, S, _, k, T, I
                      return Oe(this, function (P) {
                        switch (P.label) {
                          case 0:
                            return (
                              (r = new Date().getTime()),
                              (s = b()),
                              (u = b()),
                              0 === (d = c[0]).indexOf("http://") ||
                              0 === d.indexOf("https://")
                                ? (i = new URL(d))
                                : 0 === d.indexOf("//")
                                  ? (i = new URL(
                                      ""
                                        .concat(window.location.protocol)
                                        .concat(d),
                                    ))
                                  : ((i = new URL(
                                      window.location.href,
                                    )).pathname = d),
                              (h = !1),
                              (null === (I = e.noTraceUrls) || void 0 === I
                                ? void 0
                                : I.length) &&
                                (h = e.noTraceUrls.some(function (e) {
                                  if ("string" == typeof e) {
                                    if (e === d) return !0
                                  } else if (e instanceof RegExp && e.test(d))
                                    return !0
                                })),
                              (v = !(
                                h ||
                                ([a.ERROR, a.ERRORS, a.PERF, a.SEGMENTS].indexOf(
                                  i.pathname,
                                ) > -1 &&
                                  !e.traceSDKInternal)
                              )) &&
                                (console.log("trace:".concat(s)),
                                (g = String(be(s))),
                                (m = String(be(u))),
                                (y = String(be(o.service))),
                                (w = String(be(o.serviceInstance))),
                                (E = String(be(e.pagePath))),
                                (O = String(be(i.host))),
                                (R = o.spans.length),
                                (S = ""
                                  .concat(1, "-")
                                  .concat(g, "-")
                                  .concat(m, "-")
                                  .concat(R, "-")
                                  .concat(y, "-")
                                  .concat(w, "-")
                                  .concat(E, "-")
                                  .concat(O)),
                                c[1] || ((c[1] = {}), (c[1].headers = {})),
                                (c[1].headers.sw8 = S)),
                              [4, n.apply(void 0, c)]
                            )
                          case 1:
                            return (
                              (_ = P.sent()),
                              v &&
                                ((k = new Date().getTime()),
                                (T = {
                                  operationName: e.pagePath,
                                  startTime: r,
                                  endTime: k,
                                  spanId: o.spans.length,
                                  spanLayer: l,
                                  spanType: f,
                                  isError: 0 === _.status || _.status >= 400,
                                  parentSpanId: o.spans.length - 1,
                                  componentId: p,
                                  peer: i.host,
                                  tags: e.detailMode
                                    ? [
                                        {
                                          key: "http.method",
                                          value: c[1].method,
                                        },
                                        { key: "url", value: _.url },
                                      ]
                                    : void 0,
                                }),
                                (o = we(we({}, o), {
                                  traceId: s,
                                  traceSegmentId: u,
                                })).spans.push(T),
                                t.push(o)),
                              [2, _]
                            )
                        }
                      })
                    })
                  })(c)
                }
              })(e, r)),
              (window.onbeforeunload = function (t) {
                r.length && new y("SEGMENTS", e.collector, n).sendByXhr(r)
              }),
              setInterval(function () {
                r.length &&
                  (new y("SEGMENTS", e.collector, n).sendByXhr(r), (r = []))
              }, e.sendSegmentsTime || 3e5)
          }
          var Se = new ((function () {
              function e() {}
              return (
                (e.prototype.recordPv = function (e) {
                  var t = e.collector,
                    r = e.group,
                    n = e.subGroup,
                    o = e.service,
                    i = parseInt("".concat(1e4 * Math.random()), 10) + 3e3
                  setTimeout(function () {
                    new y("PV", t).sendByXhr({
                      group: r,
                      subGroup: n,
                      service: o,
                    })
                  }, i)
                }),
                e
              )
            })())(),
            _e = function () {
              return (
                (_e =
                  Object.assign ||
                  function (e) {
                    for (var t, r = 1, n = arguments.length; r < n; r++)
                      for (var o in (t = arguments[r]))
                        Object.prototype.hasOwnProperty.call(t, o) &&
                          (e[o] = t[o])
                    return e
                  }),
                _e.apply(this, arguments)
              )
            },
            ke = new ((function () {
              function t() {
                ;(this.registered = !1),
                  (this.customOptions = {
                    collector: "//apm.hexin.cn/skywalking-web",
                    group: "unknown",
                    subGroup: "unknown",
                    rate: 1e-4,
                    jsErrors: !0,
                    apiErrors: !0,
                    resourceErrors: !0,
                    autoTracePerf: !0,
                    useFmp: !1,
                    enableSPA: !1,
                    traceSDKInternal: !1,
                    perfEntries: !1,
                    detailMode: !0,
                    noErrorPaths: [],
                    maxReqCostTime: 6e3,
                  })
              }
              return (
                (t.prototype.register = function (t) {
                  this.registered ||
                    ((this.registered = !0),
                    console.time("完成注册"),
                    this._authConfig(t),
                    this._catchErrors(this.customOptions),
                    this.customOptions.enableSPA ||
                      (this.customOptions.rate &&
                        e() < this.customOptions.rate &&
                        ((this.customOptions.rate = 1), this._performance())),
                    Re(this.customOptions),
                    ve.handleWhiteScreen(this.customOptions),
                    this._pv(),
                    console.timeEnd("完成注册"))
                }),
                (t.prototype.setPerformance = function (e) {
                  this._authConfig(e), this._performance()
                }),
                (t.prototype._catchErrors = function (e) {
                  var t = this._getHandleOpts(e)
                  e.jsErrors &&
                    (H.handleErrors(t),
                    q.handleErrors(t),
                    e.vue && z.handleErrors(t, e.vue)),
                    e.apiErrors && G.handleError(t),
                    e.resourceErrors && X.handleErrors(t)
                }),
                (t.prototype._pv = function () {
                  e() > 0.001 || Se.recordPv(this.customOptions)
                }),
                (t.prototype._performance = function () {
                  var t = this,
                    r = this.customOptions.rate
                  if (r && e() > r) return !1
                  "complete" === document.readyState
                    ? fe.recordPerf(this.customOptions)
                    : window.addEventListener(
                        "load",
                        function () {
                          return fe.recordPerf(t.customOptions)
                        },
                        !1,
                      ),
                    this.customOptions.enableSPA &&
                      window.addEventListener(
                        "hashchange",
                        function () {
                          return fe.recordPerf(t.customOptions)
                        },
                        !1,
                      )
                }),
                (t.prototype._authConfig = function (e) {
                  var t,
                    r,
                    n = e.collector,
                    o = e.service,
                    i = e.serviceVersion,
                    a = e.group,
                    c = e.subGroup,
                    s = e.sendErrorTime
                  if (n) {
                    e.collector = n && n.replace(/\/+$/, "")
                  }
                  if (
                    (i &&
                      (e.serviceVersion = (function (e) {
                        return !e || e.indexOf("offline") > -1
                          ? e
                          : ""
                              .concat(e)
                              .concat(window._offlinePkg ? ".hxapp" : "")
                      })(i)),
                    o)
                  ) {
                    var u = o.split("::"),
                      l = u.length
                    !a && l > 1 && u[0] && (e.group = u[0]),
                      !c && l > 2 && u[1] && (e.subGroup = u[1]),
                      (e.service = u.pop())
                  }
                  "object" != typeof e.hooks && (e.hooks = {}),
                    (null === (t = e.hooks) || void 0 === t
                      ? void 0
                      : t.reportHook) &&
                      e.reportHook &&
                      (e.hooks.reportHook = e.reportHook),
                    (null === (r = e.hooks) || void 0 === r
                      ? void 0
                      : r.filterHook) &&
                      e.filterHook &&
                      (e.hooks.filterHook = e.filterHook),
                    s && s < 3e3 && (s = 3e3),
                    (this.customOptions = _e(_e({}, this.customOptions), e))
                }),
                (t.prototype._getHandleOpts = function (e) {
                  return {
                    collector: e.collector,
                    service: e.service,
                    serviceVersion: e.serviceVersion,
                    group: e.group,
                    subGroup: e.subGroup,
                    pagePath: e.pagePath,
                    sendErrorTime: e.sendErrorTime,
                    sendPerfTime: e.sendPerfTime,
                    extraMsg: e.extraMsg,
                    noErrorPaths: e.noErrorPaths,
                    hooks: e.hooks || {},
                    grade: e.grade,
                    level: e.level,
                    content: e.content,
                    maxReqCostTime: e.maxReqCostTime,
                  }
                }),
                (t.prototype.reportFrameErrors = function (e, t) {
                  this._authConfig(e)
                  var r = this._getHandleOpts(_e({}, this.customOptions))
                  ee.handleErrors(r, t)
                }),
                t
              )
            })())()
          window.ClientMonitor = ke
        },
      },
      n = {}
    function o(e) {
      var t = n[e]
      if (void 0 !== t) {
        if (void 0 !== t.error) throw t.error
        return t.exports
      }
      var i = (n[e] = { exports: {} })
      try {
        var a = { id: e, module: i, factory: r[e], require: o }
        o.i.forEach(function (e) {
          e(a)
        }),
          (i = a.module),
          a.factory.call(i.exports, i, i.exports, a.require)
      } catch (e) {
        throw ((i.error = e), e)
      }
      return i.exports
    }
    ;(o.m = r),
      (o.c = n),
      (o.i = []),
      (o.hu = function (e) {
        return e + "." + o.h() + ".hot-update.js"
      }),
      (o.hmrF = function () {
        return "main." + o.h() + ".hot-update.json"
      }),
      (o.h = function () {
        return "412591a421c4e0f487a6"
      }),
      (o.o = function (e, t) {
        return Object.prototype.hasOwnProperty.call(e, t)
      }),
      (e = {}),
      (t = "skywalking-client-js:"),
      (o.l = function (r, n, i, a) {
        if (e[r]) e[r].push(n)
        else {
          var c, s
          if (void 0 !== i)
            for (
              var u = document.getElementsByTagName("script"), l = 0;
              l < u.length;
              l++
            ) {
              var f = u[l]
              if (
                f.getAttribute("src") == r ||
                f.getAttribute("data-webpack") == t + i
              ) {
                c = f
                break
              }
            }
          c ||
            ((s = !0),
            ((c = document.createElement("script")).charset = "utf-8"),
            (c.timeout = 120),
            o.nc && c.setAttribute("nonce", o.nc),
            c.setAttribute("data-webpack", t + i),
            (c.src = r)),
            (e[r] = [n])
          var p = function (t, n) {
              ;(c.onerror = c.onload = null), clearTimeout(d)
              var o = e[r]
              if (
                (delete e[r],
                c.parentNode && c.parentNode.removeChild(c),
                o &&
                  o.forEach(function (e) {
                    return e(n)
                  }),
                t)
              )
                return t(n)
            },
            d = setTimeout(
              p.bind(null, void 0, { type: "timeout", target: c }),
              12e4,
            )
          ;(c.onerror = p.bind(null, c.onerror)),
            (c.onload = p.bind(null, c.onload)),
            s && document.head.appendChild(c)
        }
      }),
      (function () {
        var e,
          t,
          r,
          n = {},
          i = o.c,
          a = [],
          c = [],
          s = "idle",
          u = 0,
          l = []
        function f(e) {
          s = e
          for (var t = [], r = 0; r < c.length; r++) t[r] = c[r].call(null, e)
          return Promise.all(t)
        }
        function p() {
          0 == --u &&
            f("ready").then(function () {
              if (0 === u) {
                var e = l
                l = []
                for (var t = 0; t < e.length; t++) e[t]()
              }
            })
        }
        function d(e) {
          if ("idle" !== s)
            throw new Error("check() is only allowed in idle status")
          return f("check")
            .then(o.hmrM)
            .then(function (r) {
              return r
                ? f("prepare").then(function () {
                    var n = []
                    return (
                      (t = []),
                      Promise.all(
                        Object.keys(o.hmrC).reduce(function (e, i) {
                          return o.hmrC[i](r.c, r.r, r.m, e, t, n), e
                        }, []),
                      ).then(function () {
                        return (
                          (t = function () {
                            return e
                              ? v(e)
                              : f("ready").then(function () {
                                  return n
                                })
                          }),
                          0 === u
                            ? t()
                            : new Promise(function (e) {
                                l.push(function () {
                                  e(t())
                                })
                              })
                        )
                        var t
                      })
                    )
                  })
                : f(g() ? "ready" : "idle").then(function () {
                    return null
                  })
            })
        }
        function h(e) {
          return "ready" !== s
            ? Promise.resolve().then(function () {
                throw new Error(
                  "apply() is only allowed in ready status (state: " + s + ")",
                )
              })
            : v(e)
        }
        function v(e) {
          ;(e = e || {}), g()
          var n = t.map(function (t) {
            return t(e)
          })
          t = void 0
          var o = n
            .map(function (e) {
              return e.error
            })
            .filter(Boolean)
          if (o.length > 0)
            return f("abort").then(function () {
              throw o[0]
            })
          var i = f("dispose")
          n.forEach(function (e) {
            e.dispose && e.dispose()
          })
          var a,
            c = f("apply"),
            s = function (e) {
              a || (a = e)
            },
            u = []
          return (
            n.forEach(function (e) {
              if (e.apply) {
                var t = e.apply(s)
                if (t) for (var r = 0; r < t.length; r++) u.push(t[r])
              }
            }),
            Promise.all([i, c]).then(function () {
              return a
                ? f("fail").then(function () {
                    throw a
                  })
                : r
                  ? v(e).then(function (e) {
                      return (
                        u.forEach(function (t) {
                          e.indexOf(t) < 0 && e.push(t)
                        }),
                        e
                      )
                    })
                  : f("idle").then(function () {
                      return u
                    })
            })
          )
        }
        function g() {
          if (r)
            return (
              t || (t = []),
              Object.keys(o.hmrI).forEach(function (e) {
                r.forEach(function (r) {
                  o.hmrI[e](r, t)
                })
              }),
              (r = void 0),
              !0
            )
        }
        ;(o.hmrD = n),
          o.i.push(function (l) {
            var v,
              g,
              m,
              y,
              w = l.module,
              E = (function (t, r) {
                var n = i[r]
                if (!n) return t
                var o = function (o) {
                    if (n.hot.active) {
                      if (i[o]) {
                        var c = i[o].parents
                        ;-1 === c.indexOf(r) && c.push(r)
                      } else (a = [r]), (e = o)
                      ;-1 === n.children.indexOf(o) && n.children.push(o)
                    } else
                      console.warn(
                        "[HMR] unexpected require(" +
                          o +
                          ") from disposed module " +
                          r,
                      ),
                        (a = [])
                    return t(o)
                  },
                  c = function (e) {
                    return {
                      configurable: !0,
                      enumerable: !0,
                      get: function () {
                        return t[e]
                      },
                      set: function (r) {
                        t[e] = r
                      },
                    }
                  }
                for (var l in t)
                  Object.prototype.hasOwnProperty.call(t, l) &&
                    "e" !== l &&
                    Object.defineProperty(o, l, c(l))
                return (
                  (o.e = function (e) {
                    return (function (e) {
                      switch (s) {
                        case "ready":
                          f("prepare")
                        case "prepare":
                          return u++, e.then(p, p), e
                        default:
                          return e
                      }
                    })(t.e(e))
                  }),
                  o
                )
              })(l.require, l.id)
            ;(w.hot =
              ((v = l.id),
              (g = w),
              (y = {
                _acceptedDependencies: {},
                _acceptedErrorHandlers: {},
                _declinedDependencies: {},
                _selfAccepted: !1,
                _selfDeclined: !1,
                _selfInvalidated: !1,
                _disposeHandlers: [],
                _main: (m = e !== v),
                _requireSelf: function () {
                  ;(a = g.parents.slice()), (e = m ? void 0 : v), o(v)
                },
                active: !0,
                accept: function (e, t, r) {
                  if (void 0 === e) y._selfAccepted = !0
                  else if ("function" == typeof e) y._selfAccepted = e
                  else if ("object" == typeof e && null !== e)
                    for (var n = 0; n < e.length; n++)
                      (y._acceptedDependencies[e[n]] = t || function () {}),
                        (y._acceptedErrorHandlers[e[n]] = r)
                  else
                    (y._acceptedDependencies[e] = t || function () {}),
                      (y._acceptedErrorHandlers[e] = r)
                },
                decline: function (e) {
                  if (void 0 === e) y._selfDeclined = !0
                  else if ("object" == typeof e && null !== e)
                    for (var t = 0; t < e.length; t++)
                      y._declinedDependencies[e[t]] = !0
                  else y._declinedDependencies[e] = !0
                },
                dispose: function (e) {
                  y._disposeHandlers.push(e)
                },
                addDisposeHandler: function (e) {
                  y._disposeHandlers.push(e)
                },
                removeDisposeHandler: function (e) {
                  var t = y._disposeHandlers.indexOf(e)
                  t >= 0 && y._disposeHandlers.splice(t, 1)
                },
                invalidate: function () {
                  switch (((this._selfInvalidated = !0), s)) {
                    case "idle":
                      ;(t = []),
                        Object.keys(o.hmrI).forEach(function (e) {
                          o.hmrI[e](v, t)
                        }),
                        f("ready")
                      break
                    case "ready":
                      Object.keys(o.hmrI).forEach(function (e) {
                        o.hmrI[e](v, t)
                      })
                      break
                    case "prepare":
                    case "check":
                    case "dispose":
                    case "apply":
                      ;(r = r || []).push(v)
                  }
                },
                check: d,
                apply: h,
                status: function (e) {
                  if (!e) return s
                  c.push(e)
                },
                addStatusHandler: function (e) {
                  c.push(e)
                },
                removeStatusHandler: function (e) {
                  var t = c.indexOf(e)
                  t >= 0 && c.splice(t, 1)
                },
                data: n[v],
              }),
              (e = void 0),
              y)),
              (w.parents = a),
              (w.children = []),
              (a = []),
              (l.require = E)
          }),
          (o.hmrC = {}),
          (o.hmrI = {})
      })(),
      (o.p = "/"),
      (function () {
        var e,
          t,
          r,
          n,
          i,
          a = (o.hmrS_jsonp = o.hmrS_jsonp || { 179: 0 }),
          c = {}
        function s(t, r) {
          return (
            (e = r),
            new Promise(function (e, r) {
              c[t] = e
              var n = o.p + o.hu(t),
                i = new Error()
              o.l(n, function (e) {
                if (c[t]) {
                  c[t] = void 0
                  var n = e && ("load" === e.type ? "missing" : e.type),
                    o = e && e.target && e.target.src
                  ;(i.message =
                    "Loading hot update chunk " +
                    t +
                    " failed.\n(" +
                    n +
                    ": " +
                    o +
                    ")"),
                    (i.name = "ChunkLoadError"),
                    (i.type = n),
                    (i.request = o),
                    r(i)
                }
              })
            })
          )
        }
        function u(e) {
          function c(e) {
            for (
              var t = [e],
                r = {},
                n = t.map(function (e) {
                  return { chain: [e], id: e }
                });
              n.length > 0;
  
            ) {
              var i = n.pop(),
                a = i.id,
                c = i.chain,
                u = o.c[a]
              if (u && (!u.hot._selfAccepted || u.hot._selfInvalidated)) {
                if (u.hot._selfDeclined)
                  return { type: "self-declined", chain: c, moduleId: a }
                if (u.hot._main)
                  return { type: "unaccepted", chain: c, moduleId: a }
                for (var l = 0; l < u.parents.length; l++) {
                  var f = u.parents[l],
                    p = o.c[f]
                  if (p) {
                    if (p.hot._declinedDependencies[a])
                      return {
                        type: "declined",
                        chain: c.concat([f]),
                        moduleId: a,
                        parentId: f,
                      }
                    ;-1 === t.indexOf(f) &&
                      (p.hot._acceptedDependencies[a]
                        ? (r[f] || (r[f] = []), s(r[f], [a]))
                        : (delete r[f],
                          t.push(f),
                          n.push({ chain: c.concat([f]), id: f })))
                  }
                }
              }
            }
            return {
              type: "accepted",
              moduleId: e,
              outdatedModules: t,
              outdatedDependencies: r,
            }
          }
          function s(e, t) {
            for (var r = 0; r < t.length; r++) {
              var n = t[r]
              ;-1 === e.indexOf(n) && e.push(n)
            }
          }
          o.f && delete o.f.jsonpHmr, (t = void 0)
          var u = {},
            l = [],
            f = {},
            p = function (e) {
              console.warn(
                "[HMR] unexpected require(" + e.id + ") to disposed module",
              )
            }
          for (var d in r)
            if (o.o(r, d)) {
              var h,
                v = r[d],
                g = !1,
                m = !1,
                y = !1,
                w = ""
              switch (
                ((h = v ? c(d) : { type: "disposed", moduleId: d }).chain &&
                  (w = "\nUpdate propagation: " + h.chain.join(" -> ")),
                h.type)
              ) {
                case "self-declined":
                  e.onDeclined && e.onDeclined(h),
                    e.ignoreDeclined ||
                      (g = new Error(
                        "Aborted because of self decline: " + h.moduleId + w,
                      ))
                  break
                case "declined":
                  e.onDeclined && e.onDeclined(h),
                    e.ignoreDeclined ||
                      (g = new Error(
                        "Aborted because of declined dependency: " +
                          h.moduleId +
                          " in " +
                          h.parentId +
                          w,
                      ))
                  break
                case "unaccepted":
                  e.onUnaccepted && e.onUnaccepted(h),
                    e.ignoreUnaccepted ||
                      (g = new Error(
                        "Aborted because " + d + " is not accepted" + w,
                      ))
                  break
                case "accepted":
                  e.onAccepted && e.onAccepted(h), (m = !0)
                  break
                case "disposed":
                  e.onDisposed && e.onDisposed(h), (y = !0)
                  break
                default:
                  throw new Error("Unexception type " + h.type)
              }
              if (g) return { error: g }
              if (m)
                for (d in ((f[d] = v),
                s(l, h.outdatedModules),
                h.outdatedDependencies))
                  o.o(h.outdatedDependencies, d) &&
                    (u[d] || (u[d] = []), s(u[d], h.outdatedDependencies[d]))
              y && (s(l, [h.moduleId]), (f[d] = p))
            }
          r = void 0
          for (var E, O = [], b = 0; b < l.length; b++) {
            var R = l[b],
              S = o.c[R]
            S &&
              (S.hot._selfAccepted || S.hot._main) &&
              f[R] !== p &&
              !S.hot._selfInvalidated &&
              O.push({
                module: R,
                require: S.hot._requireSelf,
                errorHandler: S.hot._selfAccepted,
              })
          }
          return {
            dispose: function () {
              var e
              n.forEach(function (e) {
                delete a[e]
              }),
                (n = void 0)
              for (var t, r = l.slice(); r.length > 0; ) {
                var i = r.pop(),
                  c = o.c[i]
                if (c) {
                  var s = {},
                    f = c.hot._disposeHandlers
                  for (b = 0; b < f.length; b++) f[b].call(null, s)
                  for (
                    o.hmrD[i] = s,
                      c.hot.active = !1,
                      delete o.c[i],
                      delete u[i],
                      b = 0;
                    b < c.children.length;
                    b++
                  ) {
                    var p = o.c[c.children[b]]
                    p && (e = p.parents.indexOf(i)) >= 0 && p.parents.splice(e, 1)
                  }
                }
              }
              for (var d in u)
                if (o.o(u, d) && (c = o.c[d]))
                  for (E = u[d], b = 0; b < E.length; b++)
                    (t = E[b]),
                      (e = c.children.indexOf(t)) >= 0 && c.children.splice(e, 1)
            },
            apply: function (t) {
              for (var r in f) o.o(f, r) && (o.m[r] = f[r])
              for (var n = 0; n < i.length; n++) i[n](o)
              for (var a in u)
                if (o.o(u, a)) {
                  var c = o.c[a]
                  if (c) {
                    E = u[a]
                    for (var s = [], p = [], d = [], h = 0; h < E.length; h++) {
                      var v = E[h],
                        g = c.hot._acceptedDependencies[v],
                        m = c.hot._acceptedErrorHandlers[v]
                      if (g) {
                        if (-1 !== s.indexOf(g)) continue
                        s.push(g), p.push(m), d.push(v)
                      }
                    }
                    for (var y = 0; y < s.length; y++)
                      try {
                        s[y].call(null, E)
                      } catch (r) {
                        if ("function" == typeof p[y])
                          try {
                            p[y](r, { moduleId: a, dependencyId: d[y] })
                          } catch (n) {
                            e.onErrored &&
                              e.onErrored({
                                type: "accept-error-handler-errored",
                                moduleId: a,
                                dependencyId: d[y],
                                error: n,
                                originalError: r,
                              }),
                              e.ignoreErrored || (t(n), t(r))
                          }
                        else
                          e.onErrored &&
                            e.onErrored({
                              type: "accept-errored",
                              moduleId: a,
                              dependencyId: d[y],
                              error: r,
                            }),
                            e.ignoreErrored || t(r)
                      }
                  }
                }
              for (var w = 0; w < O.length; w++) {
                var b = O[w],
                  R = b.module
                try {
                  b.require(R)
                } catch (r) {
                  if ("function" == typeof b.errorHandler)
                    try {
                      b.errorHandler(r, { moduleId: R, module: o.c[R] })
                    } catch (n) {
                      e.onErrored &&
                        e.onErrored({
                          type: "self-accept-error-handler-errored",
                          moduleId: R,
                          error: n,
                          originalError: r,
                        }),
                        e.ignoreErrored || (t(n), t(r))
                    }
                  else
                    e.onErrored &&
                      e.onErrored({
                        type: "self-accept-errored",
                        moduleId: R,
                        error: r,
                      }),
                      e.ignoreErrored || t(r)
                }
              }
              return l
            },
          }
        }
        ;(self.webpackHotUpdateskywalking_client_js = function (t, n, a) {
          for (var s in n) o.o(n, s) && ((r[s] = n[s]), e && e.push(s))
          a && i.push(a), c[t] && (c[t](), (c[t] = void 0))
        }),
          (o.hmrI.jsonp = function (e, t) {
            r || ((r = {}), (i = []), (n = []), t.push(u)),
              o.o(r, e) || (r[e] = o.m[e])
          }),
          (o.hmrC.jsonp = function (e, c, l, f, p, d) {
            p.push(u),
              (t = {}),
              (n = c),
              (r = l.reduce(function (e, t) {
                return (e[t] = !1), e
              }, {})),
              (i = []),
              e.forEach(function (e) {
                o.o(a, e) && void 0 !== a[e]
                  ? (f.push(s(e, d)), (t[e] = !0))
                  : (t[e] = !1)
              }),
              o.f &&
                (o.f.jsonpHmr = function (e, r) {
                  t && o.o(t, e) && !t[e] && (r.push(s(e)), (t[e] = !0))
                })
          }),
          (o.hmrM = function () {
            if ("undefined" == typeof fetch)
              throw new Error("No browser support: need fetch API")
            return fetch(o.p + o.hmrF()).then(function (e) {
              if (404 !== e.status) {
                if (!e.ok)
                  throw new Error(
                    "Failed to fetch update manifest " + e.statusText,
                  )
                return e.json()
              }
            })
          })
      })()
    o("./src/index.ts")
  })()
  