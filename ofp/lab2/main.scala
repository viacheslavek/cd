import scala.collection.mutable

class Fibonacci private (val value: String) {

  def toFibonacciBinaryString(decimal: BigInt): String = {
    var a = 1
    var b = 1
    var fib = List(1)

    while (b <= decimal) {
      val c = a + b
      a = b
      b = c
      fib = c :: fib
    }
    fib = fib.tail

    val result = new mutable.StringBuilder()
    var currentNumber = decimal
    while (currentNumber > 0) {
      result.append('0')
      val currentFibonacci = fib.head
      if (currentNumber >= currentFibonacci) {
        result.update(result.length()-1, '1')
        currentNumber -= currentFibonacci
      }
      fib = fib.tail
    }

    val fibResidualLength = fib.length
    for (_ <- 1 to fibResidualLength) {
      result.append('0')
    }

    result.toString
  }

  def +(other: Fibonacci): Fibonacci = {

    val num1 = other.toInteger
    val num2 = this.toInteger
    val sumNumber = num1 + num2

    val binaryString = toFibonacciBinaryString(sumNumber)

    Fibonacci(binaryString)
  }

  def %(other: Fibonacci): Fibonacci = {
    val reversedValue1 = value.reverse
    val reversedValue2 = other.value.reverse

    val result = reversedValue1.zip(reversedValue2).map {
      case ('1', '1') => '1'
      case _ => '0'
    }.mkString

    if (result.contains('1')) Fibonacci(result.reverse)
    else Fibonacci("0")
  }


  def toInteger: BigInt = {
    var sum: BigInt = 0
    var fibPrev: BigInt = 1
    var fibCurr: BigInt = 1

    for (bit <- value.reverse) {
      if (bit == '1') {
        sum += fibCurr
      }
      val fibNext = fibPrev + fibCurr
      fibPrev = fibCurr
      fibCurr = fibNext
    }

    sum
  }

  override def toString: String = value
}

object Fibonacci {
  def apply(value: String): Fibonacci = new Fibonacci(value)
}

object Main extends App {

  // fib: 1 2 3 5 8 13 21 ...
  val num1 = Fibonacci("1") // Число 1 в фибоначчиевой системе
  val num2 = Fibonacci("10") // Число 2 в фибоначчиевой системе
  val num3 = Fibonacci("100") // Число 3 в фибоначчиевой системе
  val num4 = Fibonacci("101") // Число 4 в фибоначчиевой системе
  val num5 = Fibonacci("1000") // Число 5 в фибоначчиевой системе
  val num6 = Fibonacci("1001") // Число 6 в фибоначчиевой системе
  val num7 = Fibonacci("1010") // Число 7 в фибоначчиевой системе
  val num8 = Fibonacci("10000") // Число 8 в фибоначчиевой системе
  val num9 = Fibonacci("10001") // Число 9 в фибоначчиевой системе
  val num10 = Fibonacci("10010") // Число 10 в фибоначчиевой системе
  val num11 = Fibonacci("10100") // Число 11 в фибоначчиевой системе
  val num12 = Fibonacci("10101") // Число 12 в фибоначчиевой системе

  println("Перевод в BigInteger: " + num1.toInteger)
  println("Перевод в BigInteger: " + num2.toInteger)
  println("Перевод в BigInteger: " + num3.toInteger)
  println("Перевод в BigInteger: " + num4.toInteger)
  println("Перевод в BigInteger: " + num5.toInteger)
  println("Перевод в BigInteger: " + num6.toInteger)
  println("Перевод в BigInteger: " + num7.toInteger)
  println("Перевод в BigInteger: " + num8.toInteger)
  println("Перевод в BigInteger: " + num9.toInteger)
  println("Перевод в BigInteger: " + num10.toInteger)
  println("Перевод в BigInteger: " + num11.toInteger)
  println("Перевод в BigInteger: " + num12.toInteger)
  println()


//  Наибольшее число
  val generalNum1 = Fibonacci("101001010101")
  val generalNum2 = Fibonacci("101001000010")
  val biggerGeneralNumber1 = generalNum1 % generalNum2
  println("Наибольшее число из общих фибоначчиевых слагаемых: " + biggerGeneralNumber1)
  // Возвращает 101001000000

  val generalNum3 = Fibonacci("101001010101")
  val generalNum4 = Fibonacci("1000010")
  val biggerGeneralNumber2 = generalNum3 % generalNum4
  println("Наибольшее число из общих фибоначчиевых слагаемых: " + biggerGeneralNumber2)
  // Возвращает 1000000
  val biggerGeneralNumber3 = generalNum4 % generalNum3
  println("Наибольшее число из общих фибоначчиевых слагаемых: " + biggerGeneralNumber3)
  // Возвращает 1000000

  val generalNum5 = Fibonacci("101001010101")
  val generalNum6 = Fibonacci("1")
  val biggerGeneralNumber4 = generalNum5 % generalNum6
  println("Наибольшее число из общих фибоначчиевых слагаемых: " + biggerGeneralNumber4)
  // Возвращает 1

  val generalNum7 = Fibonacci("101001010101")
  val generalNum8 = Fibonacci("10")
  val biggerGeneralNumber5 = generalNum7 % generalNum8
  println("Наибольшее число из общих фибоначчиевых слагаемых: " + biggerGeneralNumber5)
  // Возвращает 0
  println()


//  Сумма
  val sum1_2 = num1 + num2
  println("Сумма: " + sum1_2)
  println("Перевод суммы в BigInteger: " + sum1_2.toInteger)

  val sum4_8 = num4 + num8
  println("Сумма: " + sum4_8)
  println("Перевод суммы в BigInteger: " + sum4_8.toInteger)

  val sum11_12 = num11 + num12
  println("Сумма: " + sum11_12)
  println("Перевод суммы в BigInteger: " + sum11_12.toInteger)

}