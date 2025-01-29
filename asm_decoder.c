#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

// Opens a file and returns a buffer
unsigned char *read_file(const char *filename, size_t *bytesRead) {
  // Open the file
  FILE *file;
  printf("INFO: Opening %s\n", filename);
  file = fopen(filename, "rb");
  if (file == NULL) {
    perror("Error opening file");
    fclose(file);
    return NULL;
  }
  printf("File Opened Successfully\n");

  // Get the file size by seeking to the end.
  size_t fileSize;
  fseek(file, 0, SEEK_END);
  fileSize = ftell(file);
  rewind(file);

  // allocate a buffer to store the file contents.
  unsigned char *buffer = (unsigned char *)malloc(fileSize);
  printf("INFO: Assigned buffer\n");
  if (buffer == NULL) {
    perror("Memory allocation failed\n");
    fclose(file);
    return NULL;
  }

  // Read the file contents into the buffer
  size_t bytesRead1 = fread(buffer, 1, fileSize, file);
  if (bytesRead1 != fileSize) {
    perror("Error: File reading Error from me.\n");
    free(buffer);
    fclose(file);
    return NULL;
  }
  printf("Read %zu bytes from the file\n", bytesRead1);
  *bytesRead = bytesRead1;
  fclose(file);
  return buffer;
}

// Map bits to registers.  Checks the last 3 bits as per the 8067 instruction
// manual.
char *map_register(unsigned char reg, unsigned char wbit) {
  if (wbit == 0b00000001) {
    switch (reg) {
    case 0b00000000:
      return "ax";
      break;
    case 0b00000001:
      return "cx";
      break;
    case 0b00000010:
      return "dx";
      break;
    case 0b00000011:
      return "bx";
      break;
    case 0b00000100:
      return "sp";
      break;
    case 0b00000101:
      return "bp";
      break;
    case 0b00000110:
      return "si";
      break;
    case 0b00000111:
      return "di";
      break;
    default:
      printf("Error\n");
      return NULL;
    }
  } else {
    switch (reg) {
    case 0b00000000:
      return "al";
      break;
    case 0b00000001:
      return "cl";
      break;
    case 0b00000010:
      return "dl";
      break;
    case 0b00000011:
      return "bl";
      break;
    case 0b00000100:
      return "ah";
      break;
    case 0b00000101:
      return "ch";
      break;
    case 0b00000110:
      return "dh";
      break;
    case 0b00000111:
      return "bh";
      break;
    default:
      printf("Error\n");
      return NULL;
    }
  }
}

char *map_reg_mod00(unsigned char rm) {
  switch (rm) {
  case 0b00000000:
    return "bx + si";
    break;
  case 0b00000001:
    return "bx + di";
    break;
  case 0b00000010:
    return "bp + si";
    break;
  case 0b00000011:
    return "bp + di";
    break;
  case 0b00000100:
    return "si";
    break;
  case 0b00000101:
    return "di";
    break;
  case 0b00000110:
    return "DIRECT ADDRESS";
    break;
  case 0b00000111:
    return "bx";
    break;
  default:
    printf("Error\n");
    return NULL;
  }
}

char *map_effadd(unsigned char rm) {
  switch (rm) {
  case 0b00000000:
    return "bx + si";
    break;
  case 0b00000001:
    return "bx + di";
    break;
  case 0b00000010:
    return "bp + si";
    break;
  case 0b00000011:
    return "bp + di";
    break;
  case 0b00000100:
    return "si";
    break;
  case 0b00000101:
    return "di";
    break;
  case 0b00000110:
    return "bp";
    break;
  case 0b00000111:
    return "bx";
    break;
  default:
    printf("Error\n");
    return NULL;
  }
}

int main(int argc, char *argv[]) {
  // Read byte stream from filename passed by the arguments.
  const char *filename = argv[1];

  if (argc > 2) {
    printf("Error: Too many argumnets. Only one is allowed\n");
    return 1;
  }

  size_t bytesRead = 0;

  unsigned char *buffer = read_file(filename, &bytesRead);
  if (buffer == NULL) {
    free(buffer);
    printf("Error: File read failed.");
  }

  // Example: Print the first few bytes
  for (size_t i = 0; i < bytesRead; i++) {
    // Get 4-bit op-codes
    unsigned char opcode_mask_4bit = 0b11110000;
    unsigned char opcode_4bit = buffer[i] & opcode_mask_4bit;

    switch (opcode_4bit) {
    case 0b10110000:
      unsigned char wbit_mask = 0b00001000;
      unsigned char wbit = (buffer[i] & wbit_mask) >> 3;
      unsigned char reg_mask = 0b00000111;
      unsigned char reg = (buffer[i] & reg_mask);

      switch (wbit) {
      case 0b00000000:
        printf("mov %s, %d\n", map_register(reg, wbit), buffer[i + 1]);
        break;
      case 0b00000001:
        printf("mov %s, %d\n", map_register(reg, wbit),
               (((uint16_t)buffer[i + 2] << 8) | buffer[i + 1]));
        break;
      }

    default:
      break;
    }

    // Get 6-bit op-codes
    unsigned char opcode_mask_6bit = 0b11111100;
    unsigned char opcode = buffer[i] & opcode_mask_6bit;

    // Check the op code.
    switch (opcode) {
    case 0b10001000:
      // Page 4-20 in the 8086 manual.
      // get the mod, reg and r/m bits
      unsigned char wbit_mask = 0b00000001;
      unsigned char wbit = buffer[i] & wbit_mask;
      unsigned char mod_mask = 0b11000000;
      unsigned char mod = (buffer[i + 1] & mod_mask) >> 6;
      unsigned char reg_mask = 0b00111000;
      unsigned char reg = (buffer[i + 1] & reg_mask) >> 3;
      unsigned char rm_mask = 0b00000111;
      unsigned char rm = buffer[i + 1] & rm_mask;
      unsigned char dbit_mask = 0b00000010;
      unsigned char dbit = buffer[i] & rm_mask;

      // Check the mode for the mov
      switch (mod) {
      case 0b00000011:
        printf("mov %s, %s\n", map_register(rm, wbit), map_register(reg, wbit));
        i++;
        break;
      case 0b00000001:
        printf("mov %s, [%s + %d]\n", map_register(reg, wbit), map_effadd(rm),
               buffer[i + 2]);
        i = i + 2;
        break;
      case 0b00000010:
        printf("mov %s, [%s + %d]\n", map_register(reg, wbit), map_effadd(rm),
               (((uint16_t)buffer[i + 3] << 8) | buffer[i + 2]));
        i = i + 3;
        break;
      case 0b00000000:
        printf("%b %b %b\n", buffer[i], buffer[i + 1], buffer[i + 2]);
        printf("mov %s, [%s]\n", map_register(reg, wbit), map_effadd(rm));
        i++;
        break;
      }
      break;

    default:
      break;
    }
  }
  printf("\n");

  // clean up
  free(buffer);
  return 0;
}
