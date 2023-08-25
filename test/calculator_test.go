package test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCalculator struct {
	mock.Mock
}

func (m *MockCalculator) Add(a int) int {
	args := m.Called(a)
	return args.Int(0)
}

func (m *MockCalculator) Multiply(a, b int) int {
	args := m.Called(a, b)
	return args.Int(0)
}

type Calculator interface {
	Add(a int) int
	Multiply(a, b int) int
}

func NewCalculator(cal Calculator) Calculator {
	return &CalculatorImpl{
		c: cal,
	}
}

type CalculatorImpl struct {
	c Calculator
}

func (c *CalculatorImpl) Add(a int) int {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(10)
	m := c.Multiply(5, r)
	p := a + m
	return p
}

func (c *CalculatorImpl) Multiply(a, b int) int {
	m := a * b
	return m
}

func TestMultiply(t *testing.T) {
	calculator := NewCalculator(&CalculatorImpl{})
	m := calculator.Multiply(5, 5)
	assert.Equal(t, m, 25)
}

func TestAdd(t *testing.T) {
	mockCal := new(MockCalculator)
	mockCal.On("Multiply", 2, 3).Return(6)

	Calculator := NewCalculator(mockCal)
	a := Calculator.Add(5)
	assert.Equal(t, a, 11)
	mockCal.AssertExpectations(t)
}
