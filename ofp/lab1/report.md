% Лабораторная работа № 1 «Введение в функциональное
программирование на языке Scala»
% 15 февраля 2023 г.
% Локшин Вячеслав, ИУ9-61Б

# Цель работы
Целью данной работы является ознакомление с программированием на языке Scala на основе чистых функций.

# Индивидуальный вариант
Закаренная функция

arith: (Int => Boolean) =>
(Int, Int, Int) => List[Int],
возвращающая список из n первых членов арифметической прогрессии с начальным членом a0 и разностью d,
удовлетворяющих некоторому предикату.


# Реализация и тестирование

Работа в REPL-интерпретаторе Scala:

```scala
def reverseList[A](list: List[A]): List[A] = {
  @annotation.tailrec
  def reverseHelper(src: List[A], acc: List[A]): List[A] = {
    src match {
      case Nil => acc
      case x :: xs => reverseHelper(xs, x :: acc)
    }
  }

  reverseHelper(list, Nil)
}

val arith: (Int => Boolean) => (Int, Int, Int) => List[Int] =
  pred => (n, a0, d) => {
    @annotation.tailrec
    def arithHelper(n: Int, a: Int, d: Int, result: List[Int]): List[Int] = {
      if (n <= 0) result
      else {
        if (pred(a)) arithHelper(n - 1, a + d, d, a :: result)
        else arithHelper(n, a + d, d, result)
      }
    }

    reverseList(arithHelper(n, a0, d, Nil))
  }

def main(args: Array[String]): Unit = {
  val termsArithmeticProgression = arith(_ % 2 == 0)
  val ans = termsArithmeticProgression(10, 2, 1)
  println(ans)
}
```

# Вывод
Изучил основы языка Scala и вспомнил функциональну парадигму программирования. Познакомился 
с "закарриванием" функций. В результате была написана программа, удовлетворяющая индивидуальному варианту
