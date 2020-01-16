add_url("example",
  path={
    "prefix": [
      lambda x: x + x,
    ],
    "suffix": "/?",
  },
  query={},
  tests={},
)
