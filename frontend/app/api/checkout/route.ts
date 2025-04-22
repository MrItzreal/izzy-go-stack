import { type NextRequest, NextResponse } from "next/server"
import { stripe } from "@/lib/stripe"
import { db } from "@/lib/db"
import { orders, orderItems, products } from "@/lib/db/schema"
import { auth } from "@clerk/nextjs/server"
import { createOrderSchema } from "@/lib/db/schema"
import { eq } from "drizzle-orm"

export async function POST(req: NextRequest) {
  try {
    const { userId } = auth()

    if (!userId) {
      return new NextResponse("Unauthorized", { status: 401 })
    }

    const body = await req.json()
    const validatedData = createOrderSchema.parse(body)

    // Fetch products to get their prices
    const productIds = validatedData.items.map((item) => item.productId)
    const productData = await db
      .select()
      .from(products)
      .where(productIds.map((id) => eq(products.id, id)))

    // Create a map of product id to product data
    const productMap = new Map(productData.map((product) => [product.id, product]))

    // Calculate total
    let total = 0
    for (const item of validatedData.items) {
      const product = productMap.get(item.productId)
      if (!product) {
        return new NextResponse(`Product with id ${item.productId} not found`, { status: 404 })
      }
      total += product.price * item.quantity
    }

    // Create order in database
    const [order] = await db
      .insert(orders)
      .values({
        userId: Number.parseInt(userId),
        total,
        status: "pending",
      })
      .returning()

    // Create order items
    await db.insert(orderItems).values(
      validatedData.items.map((item) => {
        const product = productMap.get(item.productId)!
        return {
          orderId: order.id,
          productId: item.productId,
          quantity: item.quantity,
          price: product.price,
        }
      }),
    )

    // Create Stripe checkout session
    const lineItems = validatedData.items.map((item) => {
      const product = productMap.get(item.productId)!
      return {
        price_data: {
          currency: "usd",
          product_data: {
            name: product.name,
            description: product.description || undefined,
          },
          unit_amount: product.price,
        },
        quantity: item.quantity,
      }
    })

    const session = await stripe.checkout.sessions.create({
      line_items: lineItems,
      mode: "payment",
      success_url: `${process.env.NEXT_PUBLIC_APP_URL}/checkout/success?session_id={CHECKOUT_SESSION_ID}`,
      cancel_url: `${process.env.NEXT_PUBLIC_APP_URL}/checkout/canceled`,
      metadata: {
        orderId: order.id.toString(),
      },
    })

    // Update order with Stripe session id
    await db.update(orders).set({ stripeSessionId: session.id }).where(eq(orders.id, order.id))

    return NextResponse.json({ url: session.url })
  } catch (error: any) {
    console.error("Checkout error:", error)
    return new NextResponse(`Checkout Error: ${error.message}`, { status: 500 })
  }
}
