
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
