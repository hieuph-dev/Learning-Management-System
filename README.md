# RESTful API for Learning Management System (LMS)

A RESTful API built with Go (Golang) and Gin framework for an online learning platform, supporting course management, user roles, payments, and progress tracking, enhanced with Redis caching, Docker deployment, and MoMo/Zalo Pay integrations.

## Key Features

- **Authentication & Authorization**: JWT-based login/register, password reset, role-based access (Admin, Instructor, Student).
- **User Management**: Profile updates, avatar upload, password management, user analytics.
- **Course Management**: CRUD operations, categorization, levels, search, ratings, and reviews.
- **Lessons**: CRUD, video lessons, ordering, previews.
- **Enrollment & Progress**: Course enrollment, progress tracking, completion certificates.
- **Payments & Orders**: Order creation, coupons, support for MoMo and Zalo Pay payment gateways, order history.
- **Coupons**: Discount types, validation rules, usage limits.
- **Analytics**: Revenue, student, course, and enrollment analytics for instructors and admins.
- **Caching**: Redis-based caching for improved API response times.
- **Containerization**: Docker support for consistent deployment and scaling.

## Technology Stack

- **Go**: High-performance language for scalable APIs.
- **Gin**: Lightweight web framework for Go.
- **PostgreSQL**: Robust relational database with GORM ORM.
- **Redis**: In-memory data store for caching.
- **Docker**: Containerization for deployment.
- **JWT**: Secure authentication with golang-jwt/jwt.
- **bcrypt**: Password hashing for security.
- **MoMo/Zalo Pay SDKs**: Payment gateway integrations for seamless transactions.

## API Documentation

Explore the full range of API endpoints and their usage in the interactive Postman documentation:

- Postman documentation: [Link to Postman documentation](https://documenter.getpostman.com/view/19784956/2sB3QJMAXQ)

## Middleware

- **Auth**: Verifies JWT and sets user context.
- **Admin/Instructor**: Role-based access checks.
- **Rate Limiter**: 5 req/s, burst 10.
- **Logger**: Logs request/response details.
- **Cache**: Redis-based caching for frequently accessed data.

## Database Models

- **User**: Info, role, status, email verification.
- **Course**: Title, pricing, metadata, stats.
- **Lesson**: Title, video, order, publish status.
- **Enrollment**: User-course relation, progress, status.
- **Order**: Transaction, payment (MoMo/Zalo Pay), coupon details.
- **Progress**: Lesson completion, watch duration.
- **Review**: Rating, comment, status.
- **Coupon**: Discount type, validation rules.

## Security

- Password hashing (bcrypt)
- JWT with expiration
- Role-based access
- Rate limiting
- SQL injection/XSS prevention
- Secure payment processing with MoMo/Zalo Pay

## Error Handling

Custom error codes (400, 401, 403, 404, 409, 500) with JSON response format.

## Performance

- Connection pooling (max 50)
- Optimized queries with GORM
- Indexed fields
- Pagination
- Rate limiting
- Redis caching for high-traffic endpoints
