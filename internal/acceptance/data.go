package acceptance

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
)

type TestData struct {
	// RandomInteger is a random integer which unique to this test case
	RandomInteger int

	// RandomString is a random 5 character string unique to this test case
	RandomString string

	// RandomID is a random UUID unique to this test case
	RandomID string

	// RandomPassword is a random password unique to this test case
	// This is not securely generated and only suitable for ephemeral test cases
	RandomPassword string

	// ResourceName is the fully qualified resource name, comprising of the
	// resource type and then the resource label
	// e.g. `azuread_application.test`
	ResourceName string

	// ResourceType is the Terraform Resource Type - `azuread_application`
	ResourceType string

	// resourceLabel is the local used for the resource - generally "test""
	resourceLabel string
}

func (t *TestData) UUID() string {
	return uuid.New().String()
}

// BuildTestData generates some test data for the given resource
func BuildTestData(t *testing.T, resourceType string, resourceLabel string) TestData {
	EnsureProvidersAreInitialised()

	testData := TestData{
		RandomInteger:  tf.AccRandTimeInt(),
		RandomString:   acctest.RandString(5),
		RandomID:       uuid.New().String(),
		RandomPassword: fmt.Sprintf("%s%s", "p@$$Wd", acctest.RandString(6)),
		ResourceName:   fmt.Sprintf("%s.%s", resourceType, resourceLabel),

		ResourceType:  resourceType,
		resourceLabel: resourceLabel,
	}

	return testData
}

// RandomIntOfLength is a random 8 to 18 digit integer which is unique to this test case
func (td *TestData) RandomIntOfLength(len int) int {
	// len should not be
	//  - greater then 18, longest a int can represent
	//  - less then 8, as that gives us YYMMDDRR
	if 8 > len || len > 18 {
		panic("Invalid Test: RandomIntOfLength: len is not between 8 or 18 inclusive")
	}

	// 18 - just return the int
	if len >= 18 {
		return td.RandomInteger
	}

	// 16-17 just strip off the last 1-2 digits
	if len >= 16 {
		return td.RandomInteger / int(math.Pow10(18-len))
	}

	// 8-15 keep len - 2 digits and add 2 characters of randomness on
	s := strconv.Itoa(td.RandomInteger)
	r := s[16:18]
	v := s[0 : len-2]
	i, _ := strconv.Atoi(v + r)

	return i
}

// RandomStringOfLength is a random 1 to 1024 character string which is unique to this test case
func (td *TestData) RandomStringOfLength(len int) string {
	// len should not be less then 1 or greater than 1024
	if 1 > len || len > 1024 {
		panic("Invalid Test: RandomStringOfLength: length argument must be between 1 and 1024 characters")
	}

	return acctest.RandString(len)
}
