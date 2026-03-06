# --- STAGE 1: BUILDER ---
# Pake Bookworm biar match sama apt-get dan glibc
FROM golang:1.26.1-bookworm AS builder

# 1. Install Tooling (tcsh wajib buat NASA)
RUN apt-get update && apt-get install -y \
    tcsh build-essential wget \
    && rm -rf /var/lib/apt/lists/*

# 2. Install CSPICE NASA
WORKDIR /opt
RUN wget https://naif.jpl.nasa.gov/pub/naif/toolkit//C/PC_Linux_GCC_64bit/packages/cspice.tar.Z \
    && zcat cspice.tar.Z | tar xf - \
    && cd cspice \
    && tcsh makeall.csh

# 3. FIX LINKER: Rename biar bisa dipanggil pake -lcspice
RUN cp /opt/cspice/lib/cspice.a /opt/cspice/lib/libcspice.a \
    && cp /opt/cspice/lib/csupport.a /opt/cspice/lib/libcsupport.a

# 4. Go Project setup
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 5. Build Go dengan flag CGO yang bener
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-I/opt/cspice/include"
ENV CGO_LDFLAGS="-L/opt/cspice/lib -lcspice -lm"

# Tambahin -ldflags buat ngecilkan ukuran binary
RUN go build -ldflags="-s -w" -o main ./cmd/api/main.go

# --- STAGE 2: RUNNER ---
FROM debian:bookworm-slim
WORKDIR /app

# Ambil binary dari builder
COPY --from=builder /app/main .

# Copy data NASA (Sangat Penting!)
COPY data/ ./data/

# Karena debian-slim minimalis, kadang butuh CA-certificates buat HTTPS
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

EXPOSE 8080

CMD ["./main"]