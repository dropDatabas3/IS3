// Jest mock for next/navigation compatible with CommonJS and named imports
let lastPushed = null;

const push = jest.fn((url) => {
  lastPushed = url;
});
const replace = jest.fn((url) => {
  lastPushed = url;
});
const back = jest.fn();
const refresh = jest.fn();
const prefetch = jest.fn();

function useRouter() {
  return { push, replace, back, refresh, prefetch };
}

function usePathname() {
  return '/';
}

function useSearchParams() {
  return new URLSearchParams('');
}

function __getLastPush() {
  return lastPushed;
}

function __reset() {
  lastPushed = null;
  push.mockClear();
  replace.mockClear();
  back.mockClear();
  refresh.mockClear();
  prefetch.mockClear();
}

module.exports = {
  useRouter,
  usePathname,
  useSearchParams,
  __getLastPush,
  __reset,
  __router: { push, replace, back, refresh, prefetch },
};
