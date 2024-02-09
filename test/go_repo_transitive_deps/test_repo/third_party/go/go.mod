module gorepotransitivedeps

go 1.21

require (
    // For this test to work this go.mod should only ever list "testify" as a dependency.
    // The point of the test is to ensure that the go.mod of "testify" (or any third-party)
    // itself is used as part of dependency resolving.
    github.com/stretchr/testify v1.8.4
)
