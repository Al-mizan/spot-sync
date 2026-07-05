
# 🚗 SpotSync – Assignment Requirements Specification

> Smart Parking & EV Charging Reservation
> 
> *A centralized platform for busy airports and malls to manage parking zones, specifically handling the high-demand reservation of limited EV charging spots.*

---

## 🛠️ Technology Stack

| Technology | Note |
| --- | --- |
| **Go (Golang)** | Version 1.22 or higher |
| **Echo** | `github.com/labstack/echo/v4` (High performance, minimalist web framework) |
| **GORM** | `gorm.io/gorm` (ORM for Go, use PostgreSQL driver) |
| **PostgreSQL** | Relational database (NeonDB, or Supabase) |
| **Validator** | `github.com/go-playground/validator/v10` (Struct validation, integrated with Echo) |
| **JWT** | `github.com/golang-jwt/jwt/v5` (Standard token generation & verification) |
| **bcrypt** | `golang.org/x/crypto/bcrypt` (Password hashing, cost 10-12) |

---


## 👥 User Roles & Permissions

| Role | Allowed Actions |
| --- | --- |
| **driver** | • Register and log in<br>• View all parking zones and availability<br>• Reserve a parking/EV spot<br>• View and cancel their own reservations |
| **admin** | • All driver permissions<br>• Create, update, and delete parking zones<br>• Set pricing for zones<br>• View all reservations in the system |

---

## 🔐 Authentication & Authorization System

- **JWT Flow:** Client sends credentials → Server validates & compares bcrypt hash → Server returns signed JWT → Client attaches token to `Authorization: Bearer <token>` header → Server middleware verifies signature & injects user data into Echo Context.
- **Security Rules:**
    - Passwords are never exposed in responses or logs.
    - Protected endpoints reject requests without a valid JWT (401 Unauthorized).
    - Role verification occurs in the Handler or Middleware before calling the Service (403 Forbidden).

---

## 🗄️ Database Schema Design

### Table 1: `users`

| Field | Requirement (Plain Text) |
| --- | --- |
| `id` | Auto-incrementing unique identifier for each account |
| `name` | Full display name of the user, must be provided |
| `email` | Valid login address, must be unique, must be provided |
| `password` | Encrypted string (bcrypt), must be provided during registration |
| `role` | Determines system access level, defaults to `driver`. Must be `driver` or `admin` |
| `created_at` | Timestamp marking when the account was created, auto-generated |
| `updated_at` | Timestamp marking when the account was last updated, auto-refreshed |

### Table 2: `parking_zones`

| Field | Requirement (Plain Text) |
| --- | --- |
| `id` | Auto-incrementing unique identifier for each zone |
| `name` | Descriptive name (e.g., "Terminal 1 EV Charging"), must be provided |
| `type` | Categorizes the zone, must be `general`, `ev_charging`, or `covered` |
| `total_capacity` | Maximum number of vehicles allowed in this zone simultaneously (integer, > 0) |
| `price_per_hour` | Cost to park in this zone (float/decimal, > 0) |
| `created_at` | Timestamp marking when the zone was created, auto-generated |
| `updated_at` | Timestamp marking when the zone was last updated, auto-refreshed |

### Table 3: `reservations`

| Field | Requirement (Plain Text) |
| --- | --- |
| `id` | Auto-incrementing unique identifier for each reservation |
| `user_id` | References the driver who made the booking (Foreign Key) |
| `zone_id` | References the parking zone booked (Foreign Key) |
| `license_plate` | The vehicle's license plate, must be provided, max 15 chars |
| `status` | Current state, defaults to `active`. Must be `active`, `completed`, or `cancelled` |
| `created_at` | Timestamp marking when the reservation was created, auto-generated |
| `updated_at` | Timestamp marking when the reservation was last updated, auto-refreshed |

---

## 🌐 API Endpoints Specification

### 🔹 Authentication Module

### 1. User Registration

**Access:** Public  
**Endpoint:** `POST /api/v1/auth/register`

**Request Body**
```json
{
  "name": "John Doe",
  "email": "john.doe@spotsync.com",
  "password": "securePassword123",
  "role": "driver"
}
```

**Success Response (201 Created)**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john.doe@spotsync.com",
    "role": "driver",
    "created_at": "2026-06-20T09:00:00Z",
    "updated_at": "2026-06-20T09:00:00Z"
  }
}
```

---

### 2. User Login

**Access:** Public  
**Endpoint:** `POST /api/v1/auth/login`

**Request Body**
```json
{
  "email": "john.doe@spotsync.com",
  "password": "securePassword123"
}
```

**Success Response (200 OK)**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@spotsync.com",
      "role": "driver"
    }
  }
}
```
> 💡 **Hint:** When signing the JWT during login, include the user's `id` and `role` in the token payload. These fields will be needed later to identify the requester and enforce permissions.

---

### 🔹 Parking Zones Module

### 3. Create Parking Zone

**Access:** Admin Only  
**Endpoint:** `POST /api/v1/zones`

**Request Body**
```json
{
  "name": "Terminal 1 EV Charging",
  "type": "ev_charging",
  "total_capacity": 20,
  "price_per_hour": 5.50
}
```

**Success Response (201 Created)**
```json
{
  "success": true,
  "message": "Parking zone created successfully",
  "data": {
    "id": 5,
    "name": "Terminal 1 EV Charging",
    "type": "ev_charging",
    "total_capacity": 20,
    "price_per_hour": 5.50,
    "created_at": "2026-06-20T10:30:00Z",
    "updated_at": "2026-06-20T10:30:00Z"
  }
}
```

---

### 4. Get All Parking Zones

**Access:** Public  
**Endpoint:** `GET /api/v1/zones`

**Success Response (200 OK)**
```json
{
  "success": true,
  "message": "Parking zones retrieved successfully",
  "data": [
    {
      "id": 5,
      "name": "Terminal 1 EV Charging",
      "type": "ev_charging",
      "total_capacity": 20,
      "available_spots": 14,
      "price_per_hour": 5.50,
      "created_at": "2026-06-20T10:30:00Z"
    }
  ]
}
```
> 💡 **Hint:** The `available_spots` field must be calculated dynamically (`total_capacity` minus the count of `active` reservations for that zone). You can achieve this in GORM using a subquery in the `Select` clause or by calculating it in the Service layer.

---

### 5. Get Single Parking Zone

**Access:** Public  
**Endpoint:** `GET /api/v1/zones/:id`

**Success Response (200 OK)**
*(Returns the same structure as a single item in the list above, including `available_spots`)*

---

### 🔹 Reservations Module (The Core Logic)

### 6. Reserve Parking Spot (⚠️ Concurrency Critical)

**Access:** Authenticated Users (`driver`, `admin`)  
**Endpoint:** `POST /api/v1/reservations`

**Request Body**
```json
{
  "zone_id": 5,
  "license_plate": "ABC-1234"
}
```

**Success Response (201 Created)**
```json
{
  "success": true,
  "message": "Reservation confirmed successfully",
  "data": {
    "id": 105,
    "user_id": 1,
    "zone_id": 5,
    "license_plate": "ABC-1234",
    "status": "active",
    "created_at": "2026-06-20T15:30:00Z",
    "updated_at": "2026-06-20T15:30:00Z"
  }
}
```

> 🚨 **CRITICAL CONCURRENCY RULE (The "EV Spot Bottleneck" Problem):**
> You must ensure a parking zone is never over-capacity. If `total_capacity` is 20, and 20 cars have active reservations, the 21st must be rejected.
> **The Race Condition:** If two drivers try to reserve the very last EV spot at the *exact same millisecond*, both requests might read "19 active" and both will succeed, resulting in 21 cars in a 20-car zone.
> **Implementation Requirement:** You **MUST** use a **GORM Database Transaction** combined with **Row-Level Locking** (`FOR UPDATE`) on the parking zone record to safely check capacity and create the reservation atomically.
>
> ```go
> // Pseudo-code hint for your Repository/Service
> db.Transaction(func(tx *gorm.DB) error {
>     var zone models.ParkingZone
>     // 1. Lock the row!
>     if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID).Error; err != nil {
>         return err
>     }
>     // 2. Count current 'active' reservations for this zone
>     // 3. Check if active_count < zone.total_capacity
>     // 4. If yes, create reservation. If no, return custom error (e.g., ErrZoneFull).
>     return nil // Commits transaction
> })
> ```

---

### 7. Get My Reservations

**Access:** Authenticated Users  
**Endpoint:** `GET /api/v1/reservations/my-reservations`

**Success Response (200 OK)**
```json
{
  "success": true,
  "message": "My reservations retrieved successfully",
  "data": [
    {
      "id": 105,
      "license_plate": "ABC-1234",
      "status": "active",
      "zone": {
        "id": 5,
        "name": "Terminal 1 EV Charging",
        "type": "ev_charging"
      },
      "created_at": "2026-06-20T15:30:00Z"
    }
  ]
}
```
> 💡 **Hint:** Use GORM `Preload("Zone")` to fetch the associated parking zone details without writing manual SQL JOINs.

---

### 8. Cancel Reservation

**Access:** Authenticated Users  
**Endpoint:** `DELETE /api/v1/reservations/:id`

**Success Response (200 OK)**
```json
{
  "success": true,
  "message": "Reservation cancelled successfully"
}
```
> 💡 **Hint:** Drivers can only cancel their *own* reservations. If they try to cancel someone else's, return 403 Forbidden. Changing the status to `cancelled` effectively frees up the spot in the zone.

---

### 9. Get All Reservations (Admin)

**Access:** Admin Only  
**Endpoint:** `GET /api/v1/reservations`

**Success Response (200 OK)**
*(Returns an array of all reservations in the system, including the `user` details and `zone` details via GORM Preloads).*

---

## 🚨 Common Response Patterns

**Standard Success Response Structure**
```json
{
  "success": true,
  "message": "Operation description",
  "data": "Response data"
}
```

**Standard Error Response Structure**
```json
{
  "success": false,
  "message": "Error description",
  "errors": "Error details"
}
```

**HTTP Status Codes**

| Code | Reason Phrase | Usage |
| --- | --- | --- |
| `200` | OK | Successful GET, PATCH, PUT, DELETE |
| `201` | Created | Successful POST (resource created) |
| `400` | Bad Request | Validation errors, invalid input, duplicate resource |
| `401` | Unauthorized | Missing, expired, or invalid JWT token |
| `403` | Forbidden | Valid token but insufficient role/permissions |
| `404` | Not Found | Requested resource does not exist |
| `409` | Conflict | Business logic conflict (e.g., Zone is full, duplicate license plate) |
| `500` | Internal Server Error | Unexpected server or database error |

---


> 💡 **Pro Tips:**
> - Use `Air` (`github.com/air-verse/air`) for hot-reloading during development.

---

Build something clean, concurrent, and well-documented.