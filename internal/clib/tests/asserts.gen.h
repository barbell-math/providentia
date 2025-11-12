#ifndef CLIB_TESTS_ASSERT
#define CLIB_TESTS_ASSERT

#include <cstring>
#include <iostream>
#include <iomanip>

#define STRINGIFY(x) #x

// Tests the supplied values are equal. If not then false will be returned.
// This macro is expected to be put in a function that returns a boolean.
#define EQ(l, r) \
	if (!Require::Eq(l, r)) { \
		Require::PrintErr( \
			__FILE__, \
			__LINE__, \
			"The supplied values were not equal but were expected to be.", \
			STRINGIFY(l), \
			STRINGIFY(r), \
			l, \
			r \
		); \
		return false; \
	}

// Tests the supplied expression(s) throw the supplied error.
// This macro is expected to be put in a function that returns a boolean.
#define THROWS(err, ...) \
	{ \
		bool wasThrown = false; \
		try { \
			__VA_ARGS__ \
		} catch (err) { wasThrown = true; } \
		catch (...) {} \
		if (!wasThrown) { \
			Require::PrintErr( \
				__FILE__, \
				__LINE__, \
				"The required exception was not thrown.", \
				"<expression>", \
				STRINGIFY(err), \
				"<no value>", \
				"<no value>" \
			); \
			return false; \
		} \
	}

namespace Require {

template <typename T, typename U>
void PrintErr(
	const char *file,
	int line,
	const char *opMsg,
	const char *lStr,
	const char *rStr,
	T l,
	U r
) {
	int lStrLen = strnlen(lStr, 100);
	int rStrLen = strnlen(rStr, 100);
	int width = lStrLen>rStrLen? lStrLen: rStrLen;

	std::cout << "\tFile: " << file << " Line: " << line << std::endl;
	std::cout << "\t" << opMsg << std::endl;
	std::cout << "\t\tLeft value : (" << std::setw(width) << lStr << ") " << l << std::endl;
	std::cout << "\t\tRight value: (" << std::setw(width) << rStr << ") " << r << std::endl;
}

template <typename T>
inline bool Eq(T l, T r) {
	return l==r;
}

};

#endif
