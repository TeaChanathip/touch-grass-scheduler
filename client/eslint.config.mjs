import js from "@eslint/js";
import { FlatCompat } from "@eslint/eslintrc";

const compat = new FlatCompat({
    baseDirectory: import.meta.dirname,
    recommendedConfig: js.configs.recommended,
});

const eslintConfig = [
    {
        ignores: [".next/**", "node_modules/**", "out/**", "build/**"],
    },
    ...compat.config({
        extends: ["eslint:recommended", "next", "prettier"],
    }),
];

export default eslintConfig;
