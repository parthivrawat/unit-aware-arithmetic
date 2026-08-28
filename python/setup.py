"""
Setup script for unit-aware-arithmetic
"""

from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="unit-aware-arithmetic",
    version="1.0.0",
    author="Parthiv Rawat",
    author_email="parthiv05022000@gmail.com",
    description="Type-safe dimensional arithmetic library with unit tracking",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/parthivrawat/unit-aware-arithmetic",
    py_modules=["dimensional"],
    classifiers=[
        "Development Status :: 5 - Production/Stable",
        "Intended Audience :: Developers",
        "Intended Audience :: Science/Research",
        "Topic :: Software Development :: Libraries :: Python Modules",
        "Topic :: Scientific/Engineering :: Physics",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
        "Typing :: Typed",
    ],
    python_requires=">=3.7",
    install_requires=[],
    extras_require={
        "dev": [
            "pytest>=7.0.0",
            "pytest-cov>=4.0.0",
            "black>=23.0.0",
            "mypy>=1.0.0",
        ],
    },
    keywords="units, dimensions, physics, measurement, type-safe, dimensional-analysis",
)
